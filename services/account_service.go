package services

import (
	"log"
	"os"
	"path/filepath"

	"main.go/config"
	"main.go/models"
)

type AccountService struct{}

func NewAccountService() *AccountService {
	return &AccountService{}
}

func (s *AccountService) DeleteAccount(email string) error {
	configuration := config.GetConfig()
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		return err
	}

	for _, filePath := range []string{user.ProfilePicture, user.DocumentImage, user.SelfieImage} {
		if filePath != "" {
			absPath := filepath.Join(configuration.UploadDirectory, filepath.Base(filePath))
			if err := os.Remove(absPath); err != nil {
				log.Printf("Failed to delete file: %v", err)
			}
		}
	}

	config.DB.Where("email = ?", email).Delete(&models.OTP{})

	return nil
}
