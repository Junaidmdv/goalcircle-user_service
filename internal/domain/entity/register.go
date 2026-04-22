package entity

import (
	"time"

	"gorm.io/gorm"
)

type TempUser struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	FullName  string
	Email     string `gorm:"unique"`
	PhoneNum  string
	Password  string
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Otp struct {
	ID         uint `gorm:"primaryKey:autoIncrement"`
	TempUserID uint `gorm:"not null"`
	Otp        string
	Type       string
	Attempts   int
	ExpiresAt  time.Time
	CreatedAt  time.Time
	DeletedAt  time.Time
	TempUser   TempUser
}
