package handlers

import (
	"chirpy/internal/config"
	"chirpy/internal/services"
	"chirpy/internal/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func NewChirpHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		chirp := services.Chirp{}
		service := services.ChirpService{Queries: cfg.Queries}
		err := decoder.Decode(&chirp)
		createChirp, err := service.CreateChirps(r.Context(), chirp.Body, chirp.UserId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("unable to create new chirp: %v", err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		err = utils.WriteJSON(w, createChirp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("unable to marshal response %v", err)
			return
		}
	}
}
func GetChirpsByID(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service := services.ChirpService{
			Queries: cfg.Queries,
		}

		id, err := uuid.Parse(r.PathValue("chirpID"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("unable to create uuid: %v", err)
			return
		}

		chirp, err := service.GetChirpId(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		err = utils.WriteJSON(w, chirp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("unable to marshal response %v", err)
			return
		}
	}
}

func GetChirpsHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service := services.ChirpService{Queries: cfg.Queries}
		chirps, err := service.GetChirps(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("unable to get chirps %v", err)
			return
		}
		err = utils.WriteJSON(w, chirps)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("unable to marshal response %v", err)
			return
		}
	}
}
