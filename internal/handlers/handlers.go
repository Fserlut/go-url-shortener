package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Fserlut/go-url-shortener/internal/config"
	"github.com/Fserlut/go-url-shortener/internal/storage"
)

type Handlers struct {
	store *storage.Storage
	cfg   *config.Config
}

type CreateShortURLRequest struct {
	URL string `json:"url"`
}

type CreateShortURLResponse struct {
	Result string `json:"result"`
}

func (h *Handlers) CreateShortURI(w http.ResponseWriter, r *http.Request) {
	// десериализуем запрос в структуру модели
	//logger.Log.Debug("decoding request")
	var req CreateShortURLRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		//logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) CreateShortURL(res http.ResponseWriter, req *http.Request) {
	//body := CreateShortURLRequest{}
	//json.Unmarshal(req.Body, body)
	body, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	url := string(body)
	if len(url) == 0 {
		res.WriteHeader(http.StatusBadRequest)
	}
	if _, ok := h.store.URLStorage[url]; ok {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	shortURL := h.store.AddURL(url)
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(fmt.Sprintf("%s/%s", h.cfg.BaseReturnURL, shortURL)))
}

func (h *Handlers) RedirectToLink(res http.ResponseWriter, req *http.Request) {
	if value, ok := h.store.URLStorage[chi.URLParam(req, "id")]; ok {
		http.Redirect(res, req, value, http.StatusTemporaryRedirect)
		return
	}
	res.WriteHeader(http.StatusBadRequest)
}

func InitHandlers(store *storage.Storage, cfg *config.Config) *Handlers {
	return &Handlers{
		cfg:   cfg,
		store: store,
	}
}
