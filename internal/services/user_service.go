package services

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
}

type UserService struct {
	Queries *database.Queries
}

func (u *UserService) CreateUser(ctx context.Context, email, password string) (User, error) {
	hash, _ := auth.HashPassword(password)
	createUser, err := u.Queries.CreateUser(ctx, database.CreateUserParams{
		Email: email,
		HashPassword: sql.NullString{
			String: hash, Valid: true,
		},
	})
	if err != nil {
		return User{}, err
	}
	return mapToUser(createUser), nil
}

func (u *UserService) GetUser(ctx context.Context, email, password string) (User, error) {

	getUser, err := u.Queries.GetUser(ctx, email)
	if err != nil {
		return User{}, err
	}
	ok, err := auth.CheckPasswordHash(password, getUser.HashPassword.String)
	if err != nil {
		return User{}, err
	}
	if !ok {
		return User{}, errors.New("incorrect email or password")
	}
	return mapToUser(getUser), nil
}

func (u *UserService) Login(ctx context.Context, email, password string) (User, error) {
	return u.GetUser(ctx, email, password)
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
