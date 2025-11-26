package services

import (
	"verify/config"
	"verify/models"

	"golang.org/x/crypto/bcrypt"
)

type SignupService struct{}

func NewSignupService() *SignupService {
	return &SignupService{}
}

func (s *SignupService) CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return config.DB.Create(user).Error
}
