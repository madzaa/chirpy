package handlers

import (
	"chirpy/internal/config"
	"chirpy/internal/services"
	"chirpy/internal/utils"
	"encoding/json"
	"log"
	"net/http"
)

func NewLoginHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		user := services.User{}
		service := services.UserService{
			Queries:     cfg.Queries,
			TokenSecret: cfg.JWTSecret,
		}
		err := decoder.Decode(&user)
		if err != nil {
			log.Printf("error unmarshalling JSON: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		loggedIn, err := service.Login(r.Context(), user.Email, user.Password, user.ExpiresIn)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = utils.WriteJSON(w, err.Error())
		} else {
			_ = utils.WriteJSON(w, loggedIn)
		}
	}
}

func NewUserHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		user := services.User{}
		service := services.UserService{
			Queries:     cfg.Queries,
			TokenSecret: cfg.JWTSecret,
		}
		err := decoder.Decode(&user)
		if err != nil {
			log.Printf("error unmarshalling JSON: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		createUser, err := service.CreateUser(r.Context(), user.Email, user.Password)
		if err != nil {
			log.Printf("error creating user: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		err = utils.WriteJSON(w, createUser)
		if err != nil {
			log.Printf("error marshaling response: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
