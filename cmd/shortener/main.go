package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/Fserlut/go-url-shortener/internal/compress"
	"github.com/Fserlut/go-url-shortener/internal/config"
	"github.com/Fserlut/go-url-shortener/internal/handlers"
	"github.com/Fserlut/go-url-shortener/internal/logger"
	"github.com/Fserlut/go-url-shortener/internal/storage"
)

func main() {
	cfg := config.InitConfig()

	if err := logger.Initialize(cfg.LogLevel); err != nil {
		panic(err)
	}

	store := storage.NewStorage(cfg)

	h := handlers.InitHandlers(store, cfg)
	r := chi.NewRouter()

	r.Use(handlers.WithLogging)

	r.Use(compress.GzipMiddleware)

	r.Post("/api/shorten", h.APICreateShortURL)
	r.Post("/api/shorten/batch", h.CreateBatchURLs)
	r.Post("/", h.CreateShortURL)

	r.Get("/{id}", h.RedirectToLink)
	r.Get("/ping", h.PingHandler)

	logger.Log.Info("Running server", zap.String("address", cfg.ServerAddress))

	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		panic(err)
	}
}
