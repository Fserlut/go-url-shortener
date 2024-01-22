package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Fserlut/go-url-shortener/internal/config"
	"github.com/Fserlut/go-url-shortener/internal/logger"
	"github.com/Fserlut/go-url-shortener/internal/storage"
	random "github.com/Fserlut/go-url-shortener/internal/utils"
)

type Handlers struct {
	store storage.Storage
	cfg   *config.Config
}

type CreateShortURLRequest struct {
	URL string `json:"url"`
}

type CreateShortURLResponse struct {
	Result string `json:"result"`
}

type CreateBatchShortenRequestItem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type CreateBatchShortenResponseItem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (h *Handlers) CreateShortURL(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	url := string(body)
	if len(url) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := h.store.SaveURL(storage.URLData{
		OriginalURL: url,
		UUID:        uuid.New().String(),
		ShortURL:    random.GetShortURL(),
	})

	if err != nil {
		if errors.Is(err, &storage.ErrURLExists{}) {
			res.Header().Set("content-type", "text/plain")
			res.WriteHeader(http.StatusConflict)
			_, err := res.Write([]byte(fmt.Sprintf("%s/%s", h.cfg.BaseReturnURL, data.ShortURL)))
			if err != nil {
				http.Error(res, "Failed to write response", http.StatusInternalServerError)
				return
			}
			return
		}
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(fmt.Sprintf("%s/%s", h.cfg.BaseReturnURL, data.ShortURL)))
}

func (h *Handlers) CreateBatchURLs(w http.ResponseWriter, r *http.Request) {
	var reqURLs []CreateBatchShortenRequestItem

	err := json.NewDecoder(r.Body).Decode(&reqURLs)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := make([]CreateBatchShortenResponseItem, 0, len(reqURLs))

	for _, reqItem := range reqURLs {
		shortData, err := h.store.SaveURL(storage.URLData{
			UUID:        uuid.New().String(),
			OriginalURL: reqItem.OriginalURL,
			ShortURL:    random.GetShortURL(),
		})

		if err != nil {
			logger.Log.Error("Error on save link: " + reqItem.OriginalURL)
			continue
		}

		res = append(res, CreateBatchShortenResponseItem{
			CorrelationID: reqItem.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", h.cfg.BaseReturnURL, shortData.ShortURL),
		})
	}

	resJSON, err := json.Marshal(res)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(resJSON)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) APICreateShortURL(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("decoding request")
	var req CreateShortURLRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if url := req.URL; url == "" {
		logger.Log.Error("cannot decode request JSON body")
		http.Error(w, "URL cant be empty", http.StatusBadRequest)
		return
	}

	shortURI := random.GetShortURL()

	// заполняем модель ответа
	resp := CreateShortURLResponse{
		Result: fmt.Sprintf("%s/%s", h.cfg.BaseReturnURL, shortURI),
	}

	respJSON, err := json.Marshal(resp)

	if err != nil {
		http.Error(w, "Marshaling response failed", http.StatusInternalServerError)
		return
	}

	_, err = h.store.SaveURL(storage.URLData{
		OriginalURL: req.URL,
		UUID:        uuid.New().String(),
		ShortURL:    shortURI,
	})

	if err != nil {
		if errors.Is(err, &storage.ErrURLExists{}) {
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusConflict)

			fmt.Println(respJSON)

			_, err = w.Write(respJSON)
			if err != nil {
				http.Error(w, "Failed to write response", http.StatusInternalServerError)
				return
			}
			return
		}
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(respJSON)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) RedirectToLink(res http.ResponseWriter, req *http.Request) {
	key := chi.URLParam(req, "id")
	value, err := h.store.GetShortURL(key)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	http.Redirect(res, req, value.OriginalURL, http.StatusTemporaryRedirect)
}

func (h *Handlers) PingHandler(res http.ResponseWriter, req *http.Request) {
	if err := h.store.Ping(); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func InitHandlers(store storage.Storage, cfg *config.Config) *Handlers {
	return &Handlers{
		cfg:   cfg,
		store: store,
	}
}
