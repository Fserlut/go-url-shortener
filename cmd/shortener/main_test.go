package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mainRoute(w, req)
	//for _, test := range tests {
	//	t.Run(test.name, func(t *testing.T) {
	//		httptest.NewRequest(htt)
	//		assert.Equal(t, test.want, mainRoute(test.value))
	//	})
	//}
}
