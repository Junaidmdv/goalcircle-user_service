package entity

import (
	"time"

	"gorm.io/gorm"
)

type TempUser struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	FullName  string
	Email     string `gorm:"unique"`
	Password  string
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Otp struct {
	ID        uint   `gorm:"primaryKey:autoIncrement"`
	Email     string `gorm:"not null;uniqueIndex:uq_otp_email_type"`
	Otp       string
	Type      string `gorm:"not null;uniqueIndex:uq_otp_email_type"`
	Attempts  int
	ExpiresAt time.Time
	CreatedAt time.Time
	DeletedAt time.Time
	UpdatedAt time.Time
}
