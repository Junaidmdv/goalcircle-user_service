package entity

import (
	"time"
)

type User struct {
	ID           string `gorm:"primaryKey"`
	Email        string
	Password     string
	GoogleAuthId string
	Role         string
	IsBlocked      bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}
