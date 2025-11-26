package services

import (
	"io"
	"mime/multipart"
	"os"

	"main.go/config"
	"main.go/models"
)

type UploadService struct{}

func NewUploadService() *UploadService {
	return &UploadService{}
}

func (s *UploadService) SaveFile(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func (s *UploadService) UpdateDocument(email, documentPath string) error {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}
	user.DocumentImage = documentPath
	return config.DB.Save(&user).Error
}
