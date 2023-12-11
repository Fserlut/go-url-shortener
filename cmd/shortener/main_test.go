package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
			name:   "#3 test redirect positive",
			code:   307,
			target: "/url1",
			request: request{
				method: http.MethodGet,
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
		{
			name:   "#5 bad request",
			code:   400,
			target: "/qwerty",
			request: request{
				method: http.MethodGet,
				body:   nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var bodyReader io.Reader
			if test.request.body != nil {
				if str, ok := test.request.body.(string); ok {
					bodyReader = strings.NewReader(str)
				}
			}
			request := httptest.NewRequest(test.request.method, test.target, bodyReader)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			if test.request.method == http.MethodGet {
				redirectToLink(w, request)
			} else {
				createShortLink(w, request)
			}

			res := w.Result()

			assert.Equal(t, test.code, res.StatusCode)

			defer res.Body.Close()
		})
	}
}
