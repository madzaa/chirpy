package middleware

import (
	"chirpy/internal/auth"
	"chirpy/internal/config"
	"context"
	"log"
	"net/http"
)

type ContextKey string

const UserIDKey ContextKey = "userID"

func IncrementHits(cfg *config.ApiConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.Hits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func CheckAuthToken(cfg *config.ApiConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer, err := auth.GetBearerToken(r.Header)
		if err != nil {
			log.Printf("unable to get bearer token: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userID, err := auth.ValidateJWT(bearer, cfg.JWTSecret)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("unable to validate jwt: %v bearer :%v", err, bearer)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
