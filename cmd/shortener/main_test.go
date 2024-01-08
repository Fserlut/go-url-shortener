package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Fserlut/go-url-shortener/internal/config"
	"github.com/Fserlut/go-url-shortener/internal/handlers"
	"github.com/Fserlut/go-url-shortener/internal/storage"
)

func TestHandlers(t *testing.T) {
	type request struct {
		method string
		body   interface{}
	}
	tests := []struct {
		name    string
		target  string
		code    int
		request request
	}{
		{
			name:   "#1 create url positive",
			target: "/",
			code:   201,
			request: request{
				method: http.MethodPost,
				body:   "https://google.com",
			},
		},
		{
			name:   "#2 create url bad request",
			code:   400,
			target: "/",
			request: request{
				method: http.MethodPost,
				body:   nil,
			},
		},
		{
			name:   "#4 test redirect bad request",
			code:   400,
			target: "/",
			request: request{
				method: http.MethodGet,
				body:   nil,
			},
		},
	}
	cfg := config.InitConfig()

	store := storage.InitStorage(cfg)

	h := handlers.InitHandlers(store, cfg)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var bodyReader io.Reader
			if test.request.body != nil {
				if str, ok := test.request.body.(string); ok {
					bodyReader = strings.NewReader(str)
				}
			}
			fmt.Println(test.target)
			request := httptest.NewRequest(test.request.method, test.target, bodyReader)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			if test.request.method == http.MethodGet {
				h.RedirectToLink(w, request)
			} else {
				h.CreateShortURL(w, request)
			}

			res := w.Result()

			assert.Equal(t, test.code, res.StatusCode)

			defer res.Body.Close()
		})
	}
}
