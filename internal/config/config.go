package config

import (
	"chirpy/internal/database"
	"database/sql"
	"log"
	"os"
	"sync/atomic"
)

type ApiConfig struct {
	Hits        atomic.Int32
	Queries     *database.Queries
	Env         string
	JWTSecret   string
	PolkaAPIKey string
}

func NewApiConfig() *ApiConfig {
	env := os.Getenv("PLATFORM")
	dbURL := os.Getenv("DB_URL")
	polkaApiKey := os.Getenv("POLKA_KEY")
	if polkaApiKey == "" {
		log.Fatal("api key is missing")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("jwt secret is missing")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	dbQueries := database.New(db)
	return &ApiConfig{Queries: dbQueries, Env: env, JWTSecret: jwtSecret, PolkaAPIKey: polkaApiKey}

}
