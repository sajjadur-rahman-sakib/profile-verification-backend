package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name           string
	Email          string `gorm:"unique"`
	Password       string
	Address        string
	ProfilePicture string
	DocumentImage  string
	SelfieImage    string
	IsVerified     bool
	Link           *string `gorm:"unique"`
}
