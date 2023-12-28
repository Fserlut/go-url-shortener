package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/Fserlut/go-url-shortener/internal/config"
	"github.com/Fserlut/go-url-shortener/internal/storage"
)

var urlStorage = make(map[string]string)

func redirectToLink(res http.ResponseWriter, req *http.Request) {
	if value, ok := urlStorage[chi.URLParam(req, "id")]; ok {
		http.Redirect(res, req, value, http.StatusTemporaryRedirect)
		return
	}
	res.WriteHeader(http.StatusBadRequest)
}

func createShortLink(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	url := string(body)
	if len(url) == 0 {
		res.WriteHeader(http.StatusBadRequest)
	}
	if _, ok := urlStorage[url]; ok {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	shortURL := getShortURL()
	urlStorage[shortURL] = url
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(fmt.Sprintf("%s/%s", "http://localhost:8080", shortURL)))
}

func getShortURL() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 8+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2:8]
}

func main() {
	cfg := config.InitConfig()

	store := storage.InitStorage()
	_ = store
	r := chi.NewRouter()
	r.Post("/", createShortLink)
	r.Get("/{id}", redirectToLink)

	fmt.Println("Running server on", cfg.ServerAddress)

	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		panic(err)
	}
}
