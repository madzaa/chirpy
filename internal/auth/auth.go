package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	params := argon2id.Params{
		Memory:      128 * 1024,
		Iterations:  4,
		Parallelism: uint8(runtime.NumCPU()),
		SaltLength:  16,
		KeyLength:   32,
	}
	hash, err := argon2id.CreateHash(password, &params)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	valid, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return valid, nil
}

func MakeJWT(userId uuid.UUID, tokenSecret string) (string, error) {
	now := jwt.NewNumericDate(time.Now().UTC())
	exp := jwt.NewNumericDate(time.Now().Add(60 * time.Minute).UTC())

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   userId.String(),
		IssuedAt:  now,
		ExpiresAt: exp,
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString([]byte(tokenSecret))
	return signed, err
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claim := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claim, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	}, jwt.WithLeeway(5*time.Second), jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return [16]byte{}, err
	}
	userId, err := token.Claims.GetSubject()
	if err != nil {
		return [16]byte{}, err
	}
	id, err := uuid.Parse(userId)
	if err != nil {
		return [16]byte{}, err
	}
	return id, nil
}

func getAuthToken(header http.Header, prefix string) (string, error) {
	authHeader := header.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("authorization header not found")
	}

	token, found := strings.CutPrefix(authHeader, prefix)
	if !found {
		return "", errors.New("authorization header must start with '" + prefix + "'")
	}

	if token == "" {
		return "", errors.New("token is empty")
	}

	return token, nil
}

func GetBearerToken(header http.Header) (string, error) {
	return getAuthToken(header, "Bearer ")
}

func GetPolkaApiKey(header http.Header) (string, error) {
	return getAuthToken(header, "APIKey ")
}

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	xcode := hex.EncodeToString(key)
	return xcode, nil
}
