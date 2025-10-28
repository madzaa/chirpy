package main

import (
	"chirpy/internal/config"
	"chirpy/internal/handlers"
	"chirpy/internal/middleware"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	apiConfig := config.NewApiConfig()

	mux := http.NewServeMux()

	mux.Handle("/app/", middleware.IncrementHits(apiConfig, filePath("./static")))
	mux.Handle("/app/logo.png", filePath("./static/assets"))

	mux.HandleFunc("GET /admin/metrics", handlers.NewMetricsHandler(apiConfig))
	mux.HandleFunc("POST /admin/reset", handlers.NewResetHandler(apiConfig))
	mux.HandleFunc("GET /api/healthz", handlers.NewHealthHandler)
	mux.HandleFunc("POST /api/chirps", handlers.NewChirpHandler(apiConfig))
	mux.HandleFunc("GET /api/chirps", handlers.GetChirpsHandler(apiConfig))
	mux.HandleFunc("GET /api/chirps/{chirpID}", handlers.GetChirpsByID(apiConfig))
	mux.HandleFunc("POST /api/users", handlers.NewUserHandler(apiConfig))

	server := http.Server{Handler: mux, Addr: ":8080"}
	log.Fatal(server.ListenAndServe())
}

func filePath(path string) http.Handler {
	if path == "" {
		path = "/"
	}
	return http.StripPrefix("/app/", http.FileServer(http.Dir(path)))
}
