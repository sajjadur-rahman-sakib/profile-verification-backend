package services

import (
	"time"

	"verify/config"
	"verify/models"

	"github.com/dgrijalva/jwt-go"
)

type TokenService struct{}

func NewTokenService() *TokenService {
	return &TokenService{}
}

func (s *TokenService) GenerateToken(user *models.User) (string, error) {
	configuration := config.GetConfig()

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(configuration.JwtSecret))
	if err != nil {
		return "", err
	}
	return signed, nil
}
