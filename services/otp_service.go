package services

import (
	"crypto/rand"
	"time"

	"verify/config"
	"verify/models"

	"github.com/labstack/echo/v4"
)

type OTPService struct{}

func NewOTPService() *OTPService {
	return &OTPService{}
}

func (s *OTPService) GenerateOTP(email string) (string, error) {
	otp := generateRandomOTPString(6)
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

func (s *OTPService) VerifyOTP(email, otp string) error {
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

func generateRandomOTPString(length int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	randomData := make([]byte, length)
	rand.Read(randomData)
	for i := range b {
		b[i] = charset[randomData[i]%byte(len(charset))]
	}
	return string(b)
}
