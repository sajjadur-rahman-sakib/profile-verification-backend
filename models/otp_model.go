package models

import (
	"time"

	"gorm.io/gorm"
)

type OTP struct {
	gorm.Model
	Email     string
	Code      string
	ExpiresAt time.Time
}
