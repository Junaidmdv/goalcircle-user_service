package entity

import (
	"time"
	"gorm.io/gorm"
)

type Organiser struct {
	ID                 string `gorm:"primaryKey"`
	UserID             string `gorm:"not null"`
	Name               string `gorm:"not null"`
	City               string
	PhoneNum           string
	Website            string
	LogoURL            string
	VerificationStatus string `gorm:"default:pending"`
	CreatedAt          gorm.DeletedAt
	UpdatedAt          time.Time
}
