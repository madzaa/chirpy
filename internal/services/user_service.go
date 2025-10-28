package services

import (
	"chirpy/internal/database"
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type UserService struct {
	User    *User
	Queries *database.Queries
}

func (u *UserService) CreateUser(ctx context.Context, email string) (User, error) {
	createUser, err := u.Queries.CreateUser(ctx, email)
	if err != nil {
		return User{}, err
	}
	return mapToUser(createUser), nil
}

func mapToUser(createUser database.User) User {
	responseUser := User{
		ID:        createUser.ID,
		CreatedAt: createUser.CreatedAt,
		UpdatedAt: createUser.UpdatedAt,
		Email:     createUser.Email,
	}
	return responseUser
}
