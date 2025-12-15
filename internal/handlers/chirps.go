package handlers

import (
	"chirpy/internal/middleware"
	"chirpy/internal/services"
	"chirpy/internal/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func NewChirpHandler(cfg *services.ChirpService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		chirp := services.Chirp{}
		err := decoder.Decode(&chirp)
		if err != nil {
			log.Printf("unable to decode new chirp: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		createChirp, err := cfg.Create(r.Context(), chirp.Body, userID)
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
func GetChirpsByID(cfg *services.ChirpService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("chirpID"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("unable to create uuid: %v", err)
			return
		}

		chirp, err := cfg.GetId(r.Context(), id)
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

func DeleteChirpByID(cfg *services.ChirpService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("chirpID"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("unable to parse chirpID: %v", err)
			return
		}
		err = cfg.DeleteChirp(r.Context(), id)
		if err != nil {
			if errors.Is(err, services.ErrUnauthorized) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func GetChirpsHandler(cfg *services.ChirpService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorID := r.URL.Query().Get("author_id")
		sorting := r.URL.Query().Get("sort")
		var chirps []services.Chirp
		var err error
		if authorID == "" {
			chirps, err = cfg.GetChirps(r.Context())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("unable to get chirps %v", err)
				return
			}
		} else {
			chirps, err = cfg.GetChirpsByUser(r.Context(), authorID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("unable to get chirps %v", err)
				return
			}
		}
		if sorting == "asc" {
			sort.Slice(chirps, func(i, j int) bool {
				return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
			})
		} else {
			sort.Slice(chirps, func(i, j int) bool {
				return chirps[j].CreatedAt.Before(chirps[i].CreatedAt)
			})
		}
		err = utils.WriteJSON(w, chirps)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("unable to marshal response %v", err)
			return
		}
	}
}
