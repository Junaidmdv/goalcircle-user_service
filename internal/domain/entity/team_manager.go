package entity

import "time"

type TeamManager struct {
	ID          string `gorm:"primaryKey"`
	UserID      string  ``
	Name        string
	ShortName   string
	LogoURL     string
	Email       string
	Phone       string
	City        string
	Status      string
	PlayerCount string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}