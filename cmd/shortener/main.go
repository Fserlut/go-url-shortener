package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
)

const host = "localhost:8080"

var urlStorage = make(map[string]string)

func redirectToLink(res http.ResponseWriter, req *http.Request) {
	if value, ok := urlStorage[req.URL.Path]; ok {
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
	res.Write([]byte(fmt.Sprintf("http://%s%s", host, shortURL)))
}

func getShortURL() string {
	return "/url" + strconv.Itoa(len(urlStorage)+1)
}

func main() {
	r := chi.NewRouter()
	r.Get(`/{id}`, redirectToLink)
	r.Post(`/`, createShortLink)
	urlStorage["/url1"] = "https://ya.ru"

	err := http.ListenAndServe(host, r)
	if err != nil {
		panic(err)
	}
}
