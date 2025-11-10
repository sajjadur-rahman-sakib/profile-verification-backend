package models

import (
	"gorm.io/gorm"
)

type Rating struct {
	gorm.Model
	RaterEmail string `gorm:"not null"`
	RatedEmail string `gorm:"not null"`
	Rating     int    `gorm:"not null;check:rating >= 1 AND rating <= 5"`
}
