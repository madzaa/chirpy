package services

import (
	"chirpy/internal/database"
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

type ChirpService struct {
	Queries *database.Queries
}

func (s *ChirpService) Create(ctx context.Context, body string, userID uuid.UUID) (Chirp, error) {
	if len(body) > 140 {
		return Chirp{}, fmt.Errorf("chirp too long")
	}
	cleanBody := s.checkProfanity(body)

	chirp, err := s.Queries.CreateChirps(ctx, database.CreateChirpsParams{
		Body:   cleanBody,
		UserID: userID,
	})
	if err != nil {
		return Chirp{}, err
	}
	return mapToChirp(chirp), nil
}

func (s *ChirpService) GetChirps(ctx context.Context) ([]Chirp, error) {
	chirps, err := s.Queries.GetChirps(ctx)
	if err != nil {
		return nil, err
	}
	var responseChirps []Chirp
	for _, chirp := range chirps {
		responseChirps = append(responseChirps, mapToChirp(chirp))

	}
	return responseChirps, nil
}

func (s *ChirpService) GetId(ctx context.Context, id uuid.UUID) (Chirp, error) {
	chirps, err := s.Queries.GetChirp(ctx, id)
	if err != nil {
		return Chirp{}, err
	}
	return mapToChirp(chirps), nil
}

func mapToChirp(createChirp database.Chirp) Chirp {
	responseChirp := Chirp{
		Id:        createChirp.ID,
		CreatedAt: createChirp.CreatedAt,
		UpdatedAt: createChirp.UpdatedAt,
		Body:      createChirp.Body,
		UserId:    createChirp.UserID,
	}
	return responseChirp
}

func (s *ChirpService) checkProfanity(chirp string) string {
	r, _ := regexp.Compile("(?i)(kerfuffle|sharbert|fornax)")
	out := r.ReplaceAllString(chirp, "****")
	return out
}
