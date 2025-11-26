package handlers

import (
	"log"
	"net/http"

	"verify/services"

	"github.com/labstack/echo/v4"
)

type PasswordHandler struct {
	passwordService *services.PasswordService
	otpService      *services.OTPService
	emailService    *services.EmailService
}

func NewPasswordHandler(
	passwordService *services.PasswordService,
	otpService *services.OTPService,
	emailService *services.EmailService,
) *PasswordHandler {
	return &PasswordHandler{
		passwordService: passwordService,
		otpService:      otpService,
		emailService:    emailService,
	}
}

func (h *PasswordHandler) ChangePassword(c echo.Context) error {
	email := c.FormValue("email")
	currentPassword := c.FormValue("current_password")
	newPassword := c.FormValue("new_password")
	err := h.passwordService.ChangePassword(email, currentPassword, newPassword)
	if err != nil {
		log.Printf("Change password error: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Password changed successfully"})
}

func (h *PasswordHandler) ForgotPassword(c echo.Context) error {
	email := c.FormValue("email")

	if err := h.passwordService.InitiateForgotPassword(email); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email not found"})
	}

	otp, err := h.otpService.GenerateOTP(email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate OTP"})
	}

	if err := h.emailService.SendOTP(email, otp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send OTP"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password reset OTP sent to email"})
}

func (h *PasswordHandler) ResetPassword(c echo.Context) error {
	email := c.FormValue("email")
	otp := c.FormValue("otp")
	newPassword := c.FormValue("new_password")

	if err := h.otpService.VerifyOTP(email, otp); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid or expired OTP"})
	}

	if err := h.passwordService.ResetPassword(email, newPassword); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password reset successfully"})
}
