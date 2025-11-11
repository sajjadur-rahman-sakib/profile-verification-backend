package services

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"main.go/config"
	"main.go/models"
)

type MessageService struct{}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (s *MessageService) SendMessage(senderEmail, receiverEmail, content string) error {
	if senderEmail == "" || receiverEmail == "" {
		return errors.New("sender and receiver are required")
	}
	if senderEmail == receiverEmail {
		return errors.New("cannot send message to yourself")
	}
	if len(content) == 0 {
		return errors.New("content cannot be empty")
	}

	var sender models.User
	if err := config.DB.Where("email = ?", senderEmail).First(&sender).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("sender not found")
		}
		return err
	}

	var receiver models.User
	if err := config.DB.Where("email = ? AND is_verified = ?", receiverEmail, true).First(&receiver).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("receiver not found or not verified")
		}
		return err
	}

	msg := models.Message{SenderEmail: senderEmail, SenderName: sender.Name, ReceiverEmail: receiverEmail, ReceiverName: receiver.Name, Content: content}
	return config.DB.Create(&msg).Error
}

func (s *MessageService) GetConversation(aEmail, bEmail string, since *time.Time, limit int) ([]models.Message, error) {
	if aEmail == "" || bEmail == "" {
		return nil, errors.New("both user emails are required")
	}
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	query := config.DB.Model(&models.Message{}).Where(
		"(sender_email = ? AND receiver_email = ?) OR (sender_email = ? AND receiver_email = ?)",
		aEmail, bEmail, bEmail, aEmail,
	)
	if since != nil {
		query = query.Where("created_at >= ?", *since)
	}

	var messages []models.Message
	if err := query.Order("created_at ASC").Limit(limit).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *MessageService) GetInbox(email string, limit int) ([]models.Message, error) {
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	var messages []models.Message
	if err := config.DB.Where("receiver_email = ?", email).Order("created_at DESC").Limit(limit).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
