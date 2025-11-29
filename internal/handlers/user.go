package handlers

import (
	"chirpy/internal/auth"
	"chirpy/internal/services"
	"chirpy/internal/utils"
	"encoding/json"
	"log"
	"net/http"
)

func NewLoginHandler(cfg *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		user := services.User{}
		err := decoder.Decode(&user)
		if err != nil {
			log.Printf("error unmarshalling JSON: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		loggedIn, err := cfg.Login(r.Context(), user.Email, user.Password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = utils.WriteJSON(w, err.Error())
		} else {
			_ = utils.WriteJSON(w, loggedIn)
		}
	}
}

func NewUserHandler(cfg *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		user := services.User{}
		err := decoder.Decode(&user)
		if err != nil {
			log.Printf("error unmarshalling JSON: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		switch r.Method {
		case "POST":
			create, err := cfg.Create(r.Context(), user.Email, user.Password)
			if err != nil {
				log.Printf("error creating user: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			err = utils.WriteJSON(w, create)
			if err != nil {
				log.Printf("error marshaling response: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		case "PUT":
			update, err := cfg.Update(r.Context(), user.Email, user.Password)
			if err != nil {
				log.Printf("error updating user: %s", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusOK)
			err = utils.WriteJSON(w, update)
			if err != nil {
				log.Printf("error marshaling response: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
}

func NewRefreshHandler(cfg *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			log.Printf("error getting token: %s, %s", refreshToken, err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		jwt, err := cfg.Refresh(r.Context(), refreshToken)
		if err != nil || jwt == "" {
			log.Printf("error getting user: %s", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		err = utils.WriteJSON(w, map[string]string{"token": jwt})
		if err != nil {
			log.Printf("unable to write json: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func NewRevokeHandler(cfg *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			log.Printf("error getting token: %s", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		err = cfg.Revoke(r.Context(), refreshToken)
		if err != nil {
			log.Printf("unable to revoke token: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
