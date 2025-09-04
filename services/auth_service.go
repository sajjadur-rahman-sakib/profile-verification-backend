package services

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"main.go/config"
	"main.go/models"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return config.DB.Create(user).Error
}

func (s *AuthService) GenerateOTP(email string) (string, error) {
	otp := generateRandomString(6)
	expiresAt := time.Now().Add(5 * time.Minute)

	otpModel := models.OTP{
		Email:     email,
		Code:      otp,
		ExpiresAt: expiresAt,
	}

	if err := config.DB.Create(&otpModel).Error; err != nil {
		return "", err
	}

	return otp, nil
}

func (s *AuthService) VerifyOTP(email, otp string) error {
	var otpModel models.OTP
	if err := config.DB.Where("email = ? AND code = ?", email, otp).First(&otpModel).Error; err != nil {
		return echo.ErrBadRequest
	}

	if time.Now().After(otpModel.ExpiresAt) {
		return echo.ErrBadRequest
	}

	config.DB.Delete(&otpModel)
	return nil
}

func (s *AuthService) UpdateDocument(email, documentPath string) error {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}
	user.DocumentImage = documentPath
	return config.DB.Save(&user).Error
}

func (s *AuthService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) UpdateUser(user *models.User) error {
	return config.DB.Save(user).Error
}

func (s *AuthService) Login(email, password string) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("email = ? AND is_verified = ?", email, true).First(&user).Error; err != nil {
		return nil, echo.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, echo.ErrUnauthorized
	}

	return &user, nil
}

func (s *AuthService) SaveFile(file *multipart.FileHeader, path string) error {
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

func (s *AuthService) DeleteAccount(email string) error {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		return err
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	for _, filePath := range []string{user.ProfilePicture, user.DocumentImage, user.SelfieImage} {
		if filePath != "" {
			absPath := filepath.Join(uploadDir, filepath.Base(filePath))
			if err := os.Remove(absPath); err != nil {
				log.Printf("Failed to delete file: %v", err)
			}
		}
	}

	config.DB.Where("email = ?", email).Delete(&models.OTP{})

	return nil
}

func (s *AuthService) ChangePassword(email, currentPassword, newPassword string) error {
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

func (s *AuthService) InitiateForgotPassword(email string) error {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return echo.ErrNotFound
	}
	return nil
}

func (s *AuthService) ResetPassword(email, newPassword string) error {
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

func generateRandomString(length int) string {
	b := make([]byte, length/2)
	rand.Read(b)
	return hex.EncodeToString(b)[:length]
}
