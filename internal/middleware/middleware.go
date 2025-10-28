package middleware

import (
	"chirpy/internal/config"
	"net/http"
)

func IncrementHits(cfg *config.ApiConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		cfg.Hits.Add(1)
		next.ServeHTTP(writer, request)
	})
}
