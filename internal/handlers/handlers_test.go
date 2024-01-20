package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Fserlut/go-url-shortener/internal/config"
	"github.com/Fserlut/go-url-shortener/internal/storage"
)

var cfg *config.Config
var store storage.Storage

func TestMain(m *testing.M) {
	cfg = config.InitConfig()       // Инициализация конфигурации
	store = storage.NewStorage(cfg) // Инициализация хранилища

	os.Exit(m.Run()) // Запускаем все тесты
}

func TestHandlers_APICreateShortURL(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name          string
		requestMethod string
		path          string
		contentType   string
		body          map[string]string
		want          want
	}{
		{
			name:          "#1 Create short link positive",
			requestMethod: http.MethodPost,
			path:          "/api/shorten",
			body:          map[string]string{"url": "https://google.com"},
			want: want{
				code:        http.StatusCreated,
				contentType: "application/json",
			},
		},
		{
			name:          "#1 Create short link Bad Request",
			requestMethod: http.MethodPost,
			path:          "/api/shorten",
			contentType:   "application/json",
			body:          map[string]string{"url": ""},
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	//defer store.File.Close()

	h := &Handlers{
		cfg:   cfg,
		store: store,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.body)
			request := httptest.NewRequest(test.requestMethod, test.path, bytes.NewBuffer(body))
			request.Header.Set("Content-Type", test.contentType)
			w := httptest.NewRecorder()
			h.APICreateShortURL(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
		})
	}
}

func TestHandlers_CreateShortURL(t *testing.T) {
	type want struct {
		code int
	}

	tests := []struct {
		name          string
		requestMethod string
		body          string
		path          string
		want          want
	}{
		{
			name:          "#1 Positive test",
			requestMethod: http.MethodPost,
			path:          "/",
			body:          "https://google.com",
			want: want{
				code: http.StatusCreated,
			},
		},
		{
			name:          "#2 Bad Request",
			requestMethod: http.MethodPost,
			path:          "/",
			body:          "",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	h := &Handlers{
		cfg:   cfg,
		store: store,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.requestMethod, test.path, strings.NewReader(test.body))
			w := httptest.NewRecorder()
			h.CreateShortURL(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
		})
	}
}
