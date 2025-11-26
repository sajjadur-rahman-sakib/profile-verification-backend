package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"main.go/services"
)

type OTPHandler struct {
	otpService   *services.OTPService
	emailService *services.EmailService
}

func NewOTPHandler(otpService *services.OTPService, emailService *services.EmailService) *OTPHandler {
	return &OTPHandler{
		otpService:   otpService,
		emailService: emailService,
	}
}

func (h *OTPHandler) ResendOTP(c echo.Context) error {
	email := c.FormValue("email")
	if email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email is required"})
	}

	otp, err := h.otpService.GenerateOTP(email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate OTP"})
	}

	if err := h.emailService.SendOTP(email, otp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send OTP"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "OTP resent successfully"})
}

func (h *OTPHandler) VerifyOTP(c echo.Context) error {
	email := c.FormValue("email")
	otp := c.FormValue("otp")

	if err := h.otpService.VerifyOTP(email, otp); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "OTP verified"})
}
