package entity

import "time"

type Admin struct {
	ID           string `gorm:"primaryKey"`
	FullName     string
	Email        string
	Password     string
	GoogleAuthId string
	UserType     string
	CreatedAt    time.Time  
	UpdatedAt    time.Time
}
