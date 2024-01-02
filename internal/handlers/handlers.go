package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Fserlut/go-url-shortener/internal/logger"
	"go.uber.org/zap"
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

func (h *Handlers) CreateShortURL(w http.ResponseWriter, r *http.Request) {
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

	result := fmt.Sprintf("%s/%s", h.cfg.BaseReturnURL, h.store.AddURL(req.URL))

	// заполняем модель ответа
	resp := CreateShortURLResponse{
		Result: result,
	}

	w.Header().Set("Content-Type", "application/json")

	// сериализуем ответ сервера
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		return
	}
	logger.Log.Debug("sending HTTP 200 response")
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
