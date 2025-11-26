package services

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"main.go/config"
	"main.go/models"
)

type PasswordService struct{}

func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

func (s *PasswordService) ChangePassword(email, currentPassword, newPassword string) error {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return echo.ErrUnauthorized
	}

	if len(newPassword) < 6 {
		return echo.ErrBadRequest
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return config.DB.Save(&user).Error
}

func (s *PasswordService) InitiateForgotPassword(email string) error {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return echo.ErrNotFound
	}
	return nil
}

func (s *PasswordService) ResetPassword(email, newPassword string) error {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return echo.ErrNotFound
	}

	if len(newPassword) < 6 {
		return echo.ErrBadRequest
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return config.DB.Save(&user).Error
}
