package services

import (
	"fmt"
	"net/smtp"

	"verify/config"
)

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) SendOTP(email, otp string) error {
	configuration := config.GetConfig()

	htmlContent := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; margin: 0; padding: 20px; }
				.container { max-width: 600px; margin: auto; background: #ffffff; padding: 20px; }
				.header { text-align: center; margin-bottom: 30px; }
				.verify-title { color: #6699cc; font-weight: bold; letter-spacing: 2px; font-size: 48px; }
				.otp-box { 
					background-color: #f8f9fa;
					padding: 20px;
					text-align: center;
					margin: 20px 0;
					border-radius: 5px;
				}
				.otp-code {
					font-size: 32px;
					letter-spacing: 5px;
					color: #000;
					font-weight: bold;
				}
				.security-note {
					background-color: #fff3cd;
					border: 1px solid #ffeeba;
					color: #856404;
					padding: 15px;
					margin-top: 20px;
					border-radius: 5px;
				}
				.time-note {
					color: #666;
					text-align: center;
					margin-top: 15px;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1 class="verify-title">VERIFY</h1>
					<h2>OTP Verification</h2>
					<p>To verify your email address, please use the following code:</p>
				</div>
				<div class="otp-box">
					<div class="otp-code">%s</div>
				</div>
				<div class="time-note">
					This OTP is valid for 5 minutes only.
				</div>
				<div class="security-note">
					<strong>🔒 Security Note:</strong> Never share this OTP with anyone. Our team will never ask for your OTP over call/message.
				</div>
			</div>
		</body>
		</html>
	`, otp)

	msg := []byte(fmt.Sprintf(
		"From: VERIFY <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: VERIFY\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n\r\n"+
			"%s",
		configuration.SmtpEmail, email, htmlContent,
	))

	auth := smtp.PlainAuth("", configuration.SmtpUsername, configuration.SmtpPassword, configuration.SmtpHost)
	addr := fmt.Sprintf("%s:%s", configuration.SmtpHost, configuration.SmtpPort)

	return smtp.SendMail(addr, auth, configuration.SmtpEmail, []string{email}, msg)
}
