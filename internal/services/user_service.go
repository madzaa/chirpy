package services

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"chirpy/internal/middleware"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type UserService struct {
	Queries     *database.Queries
	TokenSecret string
}

func (u *UserService) Update(ctx context.Context, email, password string) error {
	hash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}
	_, err = u.Queries.GetUser(ctx, email)
	if err != nil {
		return err
	}

	uid := ctx.Value(middleware.UserIDKey).(uuid.UUID)
	err = u.Queries.UpdateUsers(ctx, database.UpdateUsersParams{
		Email: email,
		HashPassword: sql.NullString{
			String: hash,
			Valid:  true,
		},
		ID: uid,
	})
	if err != nil {
		return err
	}
	return nil
}
func (u *UserService) Create(ctx context.Context, email, password string) (User, error) {
	hash, err := auth.HashPassword(password)
	if err != nil {
		return User{}, err
	}
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

func (u *UserService) Get(ctx context.Context, email, password string) (User, error) {

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
	user, err := u.Get(ctx, email, password)
	if err != nil {
		return User{}, err
	}

	jwt, err := auth.MakeJWT(user.ID, u.TokenSecret)
	if err != nil {
		return User{}, err
	}

	user.Token = jwt
	token, err := auth.MakeRefreshToken()
	if err != nil {
		return User{}, err
	}
	refreshToken, err := u.Queries.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		Token:  token,
		UserID: user.ID,
		ExpiresAt: sql.NullTime{
			Time: time.Now().Add(60 * 24 * time.Hour), Valid: true,
		},
		RevokedAt: sql.NullTime{
			Valid: false,
		},
	})
	if err != nil {
		return User{}, err
	}

	user.RefreshToken = refreshToken.Token
	return user, nil
}

func (u *UserService) Refresh(ctx context.Context, token string) (string, error) {
	tokenUser, err := u.Queries.GetUserFromRefreshToken(ctx, token)
	if err != nil {
		return "", err
	}

	user := User{
		ID:        tokenUser.ID,
		CreatedAt: tokenUser.CreatedAt,
		UpdatedAt: tokenUser.UpdatedAt,
		Email:     tokenUser.Email,
	}
	jwt, err := auth.MakeJWT(user.ID, u.TokenSecret)
	if err != nil {
		return "", err
	}
	return jwt, nil
}

func (u *UserService) Revoke(ctx context.Context, token string) error {
	err := u.Queries.RevokeRefreshToken(ctx, database.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Token: token,
	})
	if err != nil {
		return err
	}
	return nil
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
