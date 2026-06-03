package entity

import "time"

type Admin struct {
	ID           string `gorm:"primaryKey"`
	FullName     string
	Email        string
	Password     string
	CreatedAt    time.Time  
	UpdatedAt    time.Time
}


