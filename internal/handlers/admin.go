package handlers

import (
	"chirpy/internal/config"
	"fmt"
	"net/http"
)

func NewHealthHandler(writer http.ResponseWriter, r *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	_, err := writer.Write([]byte("OK"))
	if err != nil {
		return
	}
}

func NewResetHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cfg.Env != "dev" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		err := cfg.Queries.DeleteUsers(r.Context())
		if err != nil {
			return
		}
		cfg.Hits.Swap(0)
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("Reset to 0"))
		if err != nil {
			return
		}
	}
}

func NewMetricsHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := w.Write([]byte(fmt.Sprintf("<html>\n  <body>\n    <h1>Welcome, Chirpy Admin</h1>\n    <p>Chirpy has been visited %d times!</p>\n  </body>\n</html>", cfg.Hits.Load())))
		if err != nil {
			w.WriteHeader(503)
			return
		}
		w.WriteHeader(200)
	}
}
