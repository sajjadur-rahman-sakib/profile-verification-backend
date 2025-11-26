package handlers

import (
	"net/http"
	"path/filepath"
	"time"

	"verify/models"
	"verify/services"

	"github.com/labstack/echo/v4"
)

type SignupHandler struct {
	signupService *services.SignupService
	otpService    *services.OTPService
	emailService  *services.EmailService
	uploadService *services.UploadService
	tokenService  *services.TokenService
}

func NewSignupHandler(
	signupService *services.SignupService,
	otpService *services.OTPService,
	emailService *services.EmailService,
	uploadService *services.UploadService,
	tokenService *services.TokenService,
) *SignupHandler {
	return &SignupHandler{
		signupService: signupService,
		otpService:    otpService,
		emailService:  emailService,
		uploadService: uploadService,
		tokenService:  tokenService,
	}
}

func (h *SignupHandler) Signup(c echo.Context) error {
	file, err := c.FormFile("profile_picture")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Profile picture required"})
	}

	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")
	address := c.FormValue("address")

	profilePath := filepath.Join("uploads", time.Now().Format("20060102150405")+"_"+file.Filename)
	if err := h.uploadService.SaveFile(file, profilePath); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save profile picture"})
	}

	user := models.User{
		Name:           name,
		Email:          email,
		Password:       password,
		Address:        address,
		ProfilePicture: profilePath,
	}

	if err := h.signupService.CreateUser(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	otp, err := h.otpService.GenerateOTP(email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate OTP"})
	}

	if err := h.emailService.SendOTP(email, otp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send OTP"})
	}

	token, err := h.tokenService.GenerateToken(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":         "OTP sent to email",
		"token":           token,
		"email":           user.Email,
		"name":            user.Name,
		"profile_picture": user.ProfilePicture,
	})
}
