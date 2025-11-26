package handlers

import (
	"mime/multipart"
	"net/http"

	"github.com/labstack/echo/v4"
	"main.go/services"
)

type ProfileHandler struct {
	profileService *services.ProfileService
}

func NewProfileHandler(profileService *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

func (h *ProfileHandler) UpdateProfile(c echo.Context) error {
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

	if err := h.profileService.UpdateProfile(email, name, address, profilePicture, link); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Profile updated successfully"})
}

func (h *ProfileHandler) SearchProfile(c echo.Context) error {
	email := c.QueryParam("email")

	if email == "" {
		email = c.FormValue("email")
	}

	if email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email is required"})
	}

	user, err := h.profileService.SearchUserByEmail(email)
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
		"average_rating":  user.AverageRating,
	})
}
