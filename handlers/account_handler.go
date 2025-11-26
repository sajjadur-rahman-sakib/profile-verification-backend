package handlers

import (
	"log"
	"net/http"

	"verify/services"

	"github.com/labstack/echo/v4"
)

type AccountHandler struct {
	accountService *services.AccountService
}

func NewAccountHandler(accountService *services.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

func (h *AccountHandler) DeleteAccount(c echo.Context) error {
	email := c.FormValue("email")
	err := h.accountService.DeleteAccount(email)
	if err != nil {
		log.Printf("Delete account error: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Account deleted successfully"})
}
