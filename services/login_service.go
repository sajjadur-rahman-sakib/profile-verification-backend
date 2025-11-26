package services

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"main.go/config"
	"main.go/models"
)

type LoginService struct{}

func NewLoginService() *LoginService {
	return &LoginService{}
}

func (s *LoginService) Login(email, password string) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("email = ? AND is_verified = ?", email, true).First(&user).Error; err != nil {
		return nil, echo.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, echo.ErrUnauthorized
	}

	return &user, nil
}
