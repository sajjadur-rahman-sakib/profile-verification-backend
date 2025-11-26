package services

import (
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"verify/config"
	"verify/models"
)

type ProfileService struct {
	uploadService *UploadService
}

func NewProfileService(uploadService *UploadService) *ProfileService {
	return &ProfileService{uploadService: uploadService}
}

func (s *ProfileService) UpdateProfile(email string, name, address *string, profilePicture *multipart.FileHeader, link *string) error {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	if name != nil {
		user.Name = *name
	}

	if address != nil {
		user.Address = *address
	}

	if link != nil {
		user.Link = link
	}

	if profilePicture != nil {
		profilePath := filepath.Join("uploads", time.Now().Format("20060102150405")+"_"+profilePicture.Filename)
		if err := s.uploadService.SaveFile(profilePicture, profilePath); err != nil {
			return err
		}

		if user.ProfilePicture != "" {
			if err := os.Remove(user.ProfilePicture); err != nil {
				log.Printf("Failed to delete old profile picture: %v", err)
			}
		}

		user.ProfilePicture = profilePath
	}

	return config.DB.Save(&user).Error
}

func (s *ProfileService) SearchUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("email = ? AND is_verified = ?", email, true).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *ProfileService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *ProfileService) UpdateUser(user *models.User) error {
	return config.DB.Save(user).Error
}
