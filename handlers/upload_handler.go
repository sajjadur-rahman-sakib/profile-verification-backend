package handlers

import (
	"net/http"
	"path/filepath"
	"time"

	"verify/services"

	"github.com/labstack/echo/v4"
)

type UploadHandler struct {
	uploadService  *services.UploadService
	profileService *services.ProfileService
	faceService    *services.FaceService
}

func NewUploadHandler(
	uploadService *services.UploadService,
	profileService *services.ProfileService,
	faceService *services.FaceService,
) *UploadHandler {
	return &UploadHandler{
		uploadService:  uploadService,
		profileService: profileService,
		faceService:    faceService,
	}
}

func (h *UploadHandler) UploadDocument(c echo.Context) error {
	file, err := c.FormFile("document_image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Document image required"})
	}

	email := c.FormValue("email")
	documentPath := filepath.Join("uploads", time.Now().Format("20060102150405")+"_"+file.Filename)
	if err := h.uploadService.SaveFile(file, documentPath); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save document"})
	}

	if err := h.uploadService.UpdateDocument(email, documentPath); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Document uploaded"})
}

func (h *UploadHandler) UploadSelfie(c echo.Context) error {
	file, err := c.FormFile("selfie_image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Selfie image required"})
	}

	email := c.FormValue("email")
	selfiePath := filepath.Join("uploads", time.Now().Format("20060102150405")+"_"+file.Filename)
	if err := h.uploadService.SaveFile(file, selfiePath); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save selfie"})
	}

	user, err := h.profileService.GetUserByEmail(email)
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
	if err := h.profileService.UpdateUser(user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"verified": true, "message": "Account created successfully"})
}
