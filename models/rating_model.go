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

	RaterName           string `gorm:"-" json:"rater_name,omitempty"`
	RaterProfilePicture string `gorm:"-" json:"rater_profile_picture,omitempty"`
	RaterIsVerified     bool   `gorm:"-" json:"rater_is_verified,omitempty"`

	Email         string  `gorm:"-" json:"email,omitempty"`
	AverageRating float64 `gorm:"-" json:"average_rating,omitempty"`
	TotalRatings  int     `gorm:"-" json:"total_ratings,omitempty"`
}
