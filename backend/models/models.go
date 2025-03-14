package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:255;not null"`
	Email       string `gorm:"size:255;uniqueIndex;not null"`
	Password    string `gorm:"size:255"`
	Provider    string `gorm:"size:100;default:'local'"`
	ProviderID  string `gorm:"size:255;uniqueIndex"`
	AvatarURL   string `gorm:"size:512"`
	AccessToken string `gorm:"size:512"`
}
