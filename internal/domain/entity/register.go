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
	OTP       string
	ExpiresAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
