package handlers

import (
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/Fserlut/go-url-shortener/internal/logger"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size = size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func WithLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger.Log.Info(
			"",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.String("status", strconv.FormatInt(int64(responseData.status), 10)),
			zap.String("duration", strconv.Itoa(int(duration))),
			zap.String("size", strconv.Itoa(responseData.size)),
		)
	})
}
