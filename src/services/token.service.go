package services

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenService interface {
	GenerateToken(userID int) (string, error)
}

type tokenService struct {
}

func NewJWTService() TokenService {
	return &tokenService{}
}

func (s *tokenService) GenerateToken(userID int) (string, error) {
	payload := jwt.StandardClaims{
		Subject:   strconv.Itoa(int(userID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte("123"))
	if err != nil {
		return "", err
	}
	return token, nil
}
