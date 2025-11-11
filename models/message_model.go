package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	SenderName    string `gorm:"not null"`
	SenderEmail   string `gorm:"not null;index"`
	ReceiverName  string `gorm:"not null"`
	ReceiverEmail string `gorm:"not null;index"`
	Content       string `gorm:"type:text;not null"`
}
