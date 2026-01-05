package models

import (
	"gorm.io/gorm"
)

type Rating struct {
	gorm.Model
	RaterEmail string  `gorm:"not null" json:"rater_email"`
	RatedEmail string  `gorm:"not null" json:"rated_email"`
	Rating     int     `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Comment    *string `gorm:"type:text" json:"comment,omitempty"`
}
