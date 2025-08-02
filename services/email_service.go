package services

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) SendOTP(email, otp string) error {
	from := os.Getenv("FROM_EMAIL")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: Your OTP Code\r\n\r\nYour OTP is %s. It expires in 5 minutes.\r\n",
		from, email, otp,
	))

	auth := smtp.PlainAuth("", username, password, smtpHost)
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	return smtp.SendMail(addr, auth, from, []string{email}, msg)
}
