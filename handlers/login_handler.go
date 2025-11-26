package handlers

import (
	"log"
	"net/http"

	"verify/services"

	"github.com/labstack/echo/v4"
)

type LoginHandler struct {
	loginService *services.LoginService
	tokenService *services.TokenService
}

func NewLoginHandler(loginService *services.LoginService, tokenService *services.TokenService) *LoginHandler {
	return &LoginHandler{
		loginService: loginService,
		tokenService: tokenService,
	}
}

func (h *LoginHandler) Login(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	user, err := h.loginService.Login(email, password)
	if err != nil {
		log.Printf("Login error: %v", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	token, err := h.tokenService.GenerateToken(user)
	if err != nil {
		log.Printf("Token generation error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate token"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"name":            user.Name,
		"email":           user.Email,
		"address":         user.Address,
		"profile_picture": user.ProfilePicture,
		"link":            user.Link,
		"average_rating":  user.AverageRating,
		"token":           token,
	})
}
