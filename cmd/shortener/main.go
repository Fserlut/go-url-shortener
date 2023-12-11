package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const host = "localhost:8080"

var urlStorage = make(map[string]string)

func mainRoute(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		getHandler(res, req)
	} else if req.Method == http.MethodPost {
		postHandler(res, req)
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}

func getHandler(res http.ResponseWriter, req *http.Request) {
	if value, ok := urlStorage[req.URL.Path]; ok {
		fmt.Println(value)
		http.Redirect(res, req, value, http.StatusTemporaryRedirect)
		return
	}
	res.WriteHeader(http.StatusBadRequest)
}

func postHandler(res http.ResponseWriter, req *http.Request) {
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
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainRoute)
	urlStorage["/url1"] = "https://ya.ru"

	err := http.ListenAndServe(host, mux)
	if err != nil {
		panic(err)
	}
}
