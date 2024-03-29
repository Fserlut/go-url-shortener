package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	store, err := storage.NewStorage(cfg)

	if err != nil {
		logger.Log.Error("Error on init storage")
		panic(err)

	}

	h := handlers.InitHandlers(store, cfg)
	r := chi.NewRouter()

	r.Use(middleware.Compress(5,
		"application/javascript",
		"application/json",
		"text/css",
		"text/html",
		"text/plain",
		"text/xml"))
	r.Use(compress.GzipMiddleware)
	r.Use(handlers.WithLogging)

	r.Post("/api/shorten", h.CreateShortURLAPI)
	r.Post("/api/shorten/batch", h.CreateBatchURLs)
	r.Post("/", h.CreateShortURL)

	r.Get("/api/user/urls", h.GetUserURLs)
	r.Get("/{id}", h.RedirectToLink)
	r.Get("/ping", h.PingHandler)

	r.Delete("/api/user/urls", h.DeleteURLs)

	logger.Log.Info("Running server", zap.String("address", cfg.ServerAddress))

	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		panic(err)
	}
}
