package services

import (
	"errors"
	"time"

	"verify/config"
	"verify/models"

	"gorm.io/gorm"
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

type Contact struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (s *MessageService) GetContacts(userEmail string) ([]Contact, error) {
	if userEmail == "" {
		return nil, errors.New("email is required")
	}

	rows := make([]Contact, 0)
	raw := `
		SELECT u.email AS email, COALESCE(u.name, '') AS name
		FROM users u
		JOIN (
			SELECT counterpart_email, MAX(created_at) AS last_at
			FROM (
				SELECT m.receiver_email AS counterpart_email, m.created_at
				FROM messages m WHERE m.sender_email = ?
				UNION ALL
				SELECT m.sender_email AS counterpart_email, m.created_at
				FROM messages m WHERE m.receiver_email = ?
			) t
			GROUP BY counterpart_email
		) latest ON latest.counterpart_email = u.email
		ORDER BY latest.last_at DESC, u.email ASC
	`

	if err := config.DB.Raw(raw, userEmail, userEmail).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
