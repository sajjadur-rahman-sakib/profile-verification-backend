package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"main.go/config"
	"main.go/models"
	"main.go/services"
)

type MessageHandler struct {
	messageService *services.MessageService
}

func NewMessageHandler(messageService *services.MessageService) *MessageHandler {
	return &MessageHandler{messageService}
}

func (h *MessageHandler) SendMessage(c echo.Context) error {
	userIDVal := c.Get("user_id")
	if userIDVal == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var userID uint
	switch v := userIDVal.(type) {
	case float64:
		userID = uint(v)
	case int:
		userID = uint(v)
	case uint:
		userID = v
	case string:
		parsed, _ := strconv.Atoi(v)
		userID = uint(parsed)
	default:
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token subject"})
	}

	var sender models.User
	if err := config.DB.First(&sender, userID).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "sender not found"})
	}

	senderEmail := c.FormValue("sender_email")
	if senderEmail == "" {
		senderEmail = c.QueryParam("sender_email")
	}

	receiverEmail := c.FormValue("receiver_email")
	if receiverEmail == "" {
		receiverEmail = c.QueryParam("receiver_email")
	}
	content := c.FormValue("content")
	if content == "" {
		content = c.QueryParam("content")
	}

	if err := h.messageService.SendMessage(senderEmail, receiverEmail, content); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "message sent"})
}

func (h *MessageHandler) GetConversation(c echo.Context) error {
	userIDVal := c.Get("user_id")
	if userIDVal == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var userID uint
	switch v := userIDVal.(type) {
	case float64:
		userID = uint(v)
	case int:
		userID = uint(v)
	case uint:
		userID = v
	case string:
		parsed, _ := strconv.Atoi(v)
		userID = uint(parsed)
	default:
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token subject"})
	}

	var me models.User
	if err := config.DB.First(&me, userID).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found"})
	}

	senderEmail := c.FormValue("sender_email")
	receiverEmail := c.FormValue("receiver_email")

	if senderEmail == "" || receiverEmail == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "sender_email and receiver_email are required"})
	}

	if me.Email != senderEmail && me.Email != receiverEmail {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "not allowed to view this conversation"})
	}

	limit := 100
	if l := c.FormValue("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	var sincePtr *time.Time
	if s := c.FormValue("since"); s != "" {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			sincePtr = &t
		}
	}

	msgs, err := h.messageService.GetConversation(senderEmail, receiverEmail, sincePtr, limit)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	var sender models.User
	var receiver models.User
	if err := config.DB.Where("email = ?", senderEmail).First(&sender).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "sender not found"})
	}
	if err := config.DB.Where("email = ?", receiverEmail).First(&receiver).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "receiver not found"})
	}

	type messageDTO struct {
		Index         int    `json:"index"`
		SenderEmail   string `json:"sender_email"`
		SenderName    string `json:"sender_name"`
		ReceiverEmail string `json:"receiver_email"`
		ReceiverName  string `json:"receiver_name"`
		Content       string `json:"content"`
		CreatedAt     string `json:"created_at"`
	}

	out := make([]messageDTO, 0, len(msgs))
	for i, m := range msgs {
		out = append(out, messageDTO{
			Index:         i + 1,
			SenderEmail:   m.SenderEmail,
			SenderName:    m.SenderName,
			ReceiverEmail: m.ReceiverEmail,
			ReceiverName:  m.ReceiverName,
			Content:       m.Content,
			CreatedAt:     m.CreatedAt.Format(time.RFC3339),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"sender": map[string]interface{}{
			"email": sender.Email,
			"name":  sender.Name,
		},
		"receiver": map[string]interface{}{
			"email": receiver.Email,
			"name":  receiver.Name,
		},
		"count":    len(out),
		"messages": out,
	})
}

func (h *MessageHandler) GetContacts(c echo.Context) error {
	userIDVal := c.Get("user_id")
	if userIDVal == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var userID uint
	switch v := userIDVal.(type) {
	case float64:
		userID = uint(v)
	case int:
		userID = uint(v)
	case uint:
		userID = v
	case string:
		parsed, _ := strconv.Atoi(v)
		userID = uint(parsed)
	default:
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token subject"})
	}

	var me models.User
	if err := config.DB.First(&me, userID).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found"})
	}

	email := c.FormValue("email")
	if email == "" {
		email = c.QueryParam("email")
	}
	if email == "" {
		email = me.Email
	}

	if me.Email != email {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "not allowed to view contacts of another user"})
	}

	contacts, err := h.messageService.GetContacts(email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"email":    email,
		"count":    len(contacts),
		"contacts": contacts,
	})
}
