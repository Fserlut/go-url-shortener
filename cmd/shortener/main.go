package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Fserlut/go-url-shortener/internal/config"
	"github.com/Fserlut/go-url-shortener/internal/handlers"
	"github.com/Fserlut/go-url-shortener/internal/logger"
	"github.com/Fserlut/go-url-shortener/internal/storage"
)

func main() {
	cfg := config.InitConfig()

	log := logger.InitLogger()

	store := storage.InitStorage()

	h := handlers.InitHandlers(store, cfg)
	r := chi.NewRouter()

	r.Use(func(handler http.Handler) http.Handler {
		return handlers.WithLogging(log, handler)
	})

	r.Post("/", h.CreateShortURL)
	r.Get("/{id}", h.RedirectToLink)

	log.Infoln("Running server on", cfg.ServerAddress)

	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		panic(err)
	}
}
