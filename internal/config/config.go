package config

import (
	"chirpy/internal/database"
	"database/sql"
	"log"
	"os"
	"sync/atomic"
)

type ApiConfig struct {
	Hits      atomic.Int32
	Queries   *database.Queries
	Env       string
	JWTSecret string
}

func NewApiConfig() *ApiConfig {
	env := os.Getenv("PLATFORM")
	dbURL := os.Getenv("DB_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("jwt secret is missing")

	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)
	return &ApiConfig{Queries: dbQueries, Env: env, JWTSecret: jwtSecret}
}
