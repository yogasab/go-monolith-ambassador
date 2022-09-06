package services

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/yogasab/go-monolith-ambassador/src/middlewares"
)

type TokenService interface {
	GenerateToken(userID int, scope string) (string, error)
}

type tokenService struct {
}

func NewJWTService() TokenService {
	return &tokenService{}
}

func (s *tokenService) GenerateToken(userID int, scope string) (string, error) {
	// payload := jwt.StandardClaims{
	// 	Subject:   strconv.Itoa(int(userID)),
	// 	ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	// }
	// Custome JWT Claims
	payload := middlewares.ClaimsWithScope{}
	payload.Subject = strconv.Itoa(int(userID))
	payload.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()
	payload.Scope = scope
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte("123"))
	if err != nil {
		return "", err
	}
	return token, nil
}
