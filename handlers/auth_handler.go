package handlers

import (
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"main.go/models"
	"main.go/services"
)

type AuthHandler struct {
	authService  *services.AuthService
	emailService *services.EmailService
	faceService  *services.FaceService
}

func NewAuthHandler(authService *services.AuthService, emailService *services.EmailService, faceService *services.FaceService) *AuthHandler {
	return &AuthHandler{authService, emailService, faceService}
}

func (h *AuthHandler) Signup(c echo.Context) error {
	file, err := c.FormFile("profile_picture")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Profile picture required"})
	}

	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")
	address := c.FormValue("address")

	profilePath := filepath.Join("uploads", time.Now().Format("20060102150405")+"_"+file.Filename)
	if err := h.authService.SaveFile(file, profilePath); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save profile picture"})
	}

	user := models.User{
		Name:           name,
		Email:          email,
		Password:       password,
		Address:        address,
		ProfilePicture: profilePath,
	}

	if err := h.authService.CreateUser(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	otp, err := h.authService.GenerateOTP(email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate OTP"})
	}

	if err := h.emailService.SendOTP(email, otp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send OTP"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "OTP sent to email"})
}

func (h *AuthHandler) ResendOTP(c echo.Context) error {
	email := c.FormValue("email")
	if email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email is required"})
	}

	otp, err := h.authService.GenerateOTP(email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate OTP"})
	}

	if err := h.emailService.SendOTP(email, otp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send OTP"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "OTP resent successfully"})
}

func (h *AuthHandler) VerifyOTP(c echo.Context) error {
	email := c.FormValue("email")
	otp := c.FormValue("otp")

	if err := h.authService.VerifyOTP(email, otp); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "OTP verified"})
}

func (h *AuthHandler) UploadDocument(c echo.Context) error {
	file, err := c.FormFile("document_image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Document image required"})
	}

	email := c.FormValue("email")
	documentPath := filepath.Join("uploads", time.Now().Format("20060102150405")+"_"+file.Filename)
	if err := h.authService.SaveFile(file, documentPath); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save document"})
	}

	if err := h.authService.UpdateDocument(email, documentPath); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Document uploaded"})
}

func (h *AuthHandler) UploadSelfie(c echo.Context) error {
	file, err := c.FormFile("selfie_image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Selfie image required"})
	}

	email := c.FormValue("email")
	selfiePath := filepath.Join("uploads", time.Now().Format("20060102150405")+"_"+file.Filename)
	if err := h.authService.SaveFile(file, selfiePath); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save selfie"})
	}

	user, err := h.authService.GetUserByEmail(email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	isMatch, err := h.faceService.CompareFaces(user.DocumentImage, selfiePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to compare faces"})
	}

	if !isMatch {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"verified": false, "message": "Faces do not match"})
	}

	user.SelfieImage = selfiePath
	user.IsVerified = true
	if err := h.authService.UpdateUser(user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"verified": true, "message": "Account created successfully"})
}

func (h *AuthHandler) Login(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	user, err := h.authService.Login(email, password)
	if err != nil {
		log.Printf("Login error: %v", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"name":            user.Name,
		"email":           user.Email,
		"address":         user.Address,
		"profile_picture": user.ProfilePicture,
		"link":            user.Link,
	})
}

func (h *AuthHandler) DeleteAccount(c echo.Context) error {
	email := c.FormValue("email")
	err := h.authService.DeleteAccount(email)
	if err != nil {
		log.Printf("Delete account error: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Account deleted successfully"})
}

func (h *AuthHandler) ChangePassword(c echo.Context) error {
	email := c.FormValue("email")
	currentPassword := c.FormValue("current_password")
	newPassword := c.FormValue("new_password")
	err := h.authService.ChangePassword(email, currentPassword, newPassword)
	if err != nil {
		log.Printf("Change password error: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Password changed successfully"})
}

func (h *AuthHandler) ForgotPassword(c echo.Context) error {
	email := c.FormValue("email")

	if err := h.authService.InitiateForgotPassword(email); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email not found"})
	}

	otp, err := h.authService.GenerateOTP(email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate OTP"})
	}

	if err := h.emailService.SendOTP(email, otp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send OTP"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password reset OTP sent to email"})
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {
	email := c.FormValue("email")
	otp := c.FormValue("otp")
	newPassword := c.FormValue("new_password")

	if err := h.authService.VerifyOTP(email, otp); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid or expired OTP"})
	}

	if err := h.authService.ResetPassword(email, newPassword); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password reset successfully"})
}

func (h *AuthHandler) UpdateProfile(c echo.Context) error {
	email := c.FormValue("email")

	var name, address, link *string
	if nameStr := c.FormValue("name"); nameStr != "" {
		name = &nameStr
	}
	if addressStr := c.FormValue("address"); addressStr != "" {
		address = &addressStr
	}
	if linkStr := c.FormValue("link"); linkStr != "" {
		link = &linkStr
	}

	var profilePicture *multipart.FileHeader
	if file, err := c.FormFile("profile_picture"); err == nil {
		profilePicture = file
	}

	if name == nil && address == nil && profilePicture == nil && link == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No update data provided"})
	}

	if err := h.authService.UpdateProfile(email, name, address, profilePicture, link); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Profile updated successfully"})
}

func (h *AuthHandler) SearchProfile(c echo.Context) error {
	email := c.QueryParam("email")

	if email == "" {
		email = c.FormValue("email")
	}

	if email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email is required"})
	}

	user, err := h.authService.SearchUserByEmail(email)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Verified user not found"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"name":            user.Name,
		"email":           user.Email,
		"address":         user.Address,
		"profile_picture": user.ProfilePicture,
		"is_verified":     true,
		"link":            user.Link,
	})
}
