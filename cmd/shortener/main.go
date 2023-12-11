package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
)

var urlStorage = make(map[string]string)

func redirectToLink(res http.ResponseWriter, req *http.Request) {
	if value, ok := urlStorage[chi.URLParam(req, "id")]; ok {
		fmt.Println(value)
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
	for _, s2 := range urlStorage {
		if s2 == url {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	shortURL := getShortURL()
	urlStorage[shortURL] = url
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(fmt.Sprintf("%s/%s", baseReturnURL, shortURL)))
}

func getShortURL() string {
	return "url" + strconv.Itoa(len(urlStorage)+1)
}

func main() {
	parseFlags()
	r := chi.NewRouter()
	r.Post("/", createShortLink)
	r.Get("/{id}", redirectToLink)
	urlStorage["url1"] = "https://ya.ru"

	fmt.Println("Running server on", serverAddress)

	if err := http.ListenAndServe(serverAddress, r); err != nil {
		panic(err)
	}
}
