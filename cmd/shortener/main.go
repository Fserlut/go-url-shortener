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
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	url := string(body)
	for _, s2 := range urlStorage {
		if s2 == url {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	shortUrl := getShortUrl()
	urlStorage[shortUrl] = url
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(fmt.Sprintf("http://%s%s", host, shortUrl)))
}

func getShortUrl() string {
	return "/url" + strconv.Itoa(len(urlStorage)+1)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainRoute)

	err := http.ListenAndServe(host, mux)
	if err != nil {
		panic(err)
	}
}
