package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint64 `gorm:"primaryKey"`
	Account       string `gorm:"size:255;not null;uniqueIndex"`
	PasswordHash  string `gorm:"size:255;not null"`
	Nickname      string `gorm:"size:100;not null"`
	AuthTokenHash string `gorm:"size:64;not null"`
	State         uint8  `gorm:"not null;default:1"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
