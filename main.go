package main

import (
	"chirpy/internal/config"
	"chirpy/internal/handlers"
	"chirpy/internal/middleware"
	"chirpy/internal/services"
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
	userService := &services.UserService{
		Queries:     apiConfig.Queries,
		TokenSecret: apiConfig.JWTSecret,
	}

	chirpService := &services.ChirpService{
		Queries: apiConfig.Queries,
	}

	mux := http.NewServeMux()

	// server
	mux.Handle("GET /app/", middleware.IncrementHits(apiConfig, filePath("./static")))
	mux.Handle("GET /app/logo.png", filePath("./static/assets"))

	// chirps
	mux.HandleFunc("GET /api/chirps", handlers.GetChirpsHandler(chirpService))
	mux.HandleFunc("GET /api/chirps/{chirpID}", handlers.GetChirpsByID(chirpService))
	mux.Handle("POST /api/chirps", middleware.CheckAuthToken(apiConfig, handlers.NewChirpHandler(chirpService)))

	// admin
	mux.HandleFunc("GET /admin/metrics", handlers.NewMetricsHandler(apiConfig))
	mux.HandleFunc("POST /admin/reset", handlers.NewResetHandler(apiConfig))
	mux.HandleFunc("GET /api/healthz", handlers.NewHealthHandler)

	// users
	mux.HandleFunc("POST /api/users", handlers.NewUserHandler(userService))
	mux.HandleFunc("POST /api/login", handlers.NewLoginHandler(userService))
	mux.HandleFunc("POST /api/refresh", handlers.NewRefreshHandler(userService))
	mux.HandleFunc("POST /api/revoke", handlers.NewRevokeHandler(userService))

	server := http.Server{Handler: mux, Addr: ":8080"}
	log.Fatal(server.ListenAndServe())
}

func filePath(path string) http.Handler {
	if path == "" {
		path = "/"
	}
	return http.StripPrefix("/app/", http.FileServer(http.Dir(path)))
}
