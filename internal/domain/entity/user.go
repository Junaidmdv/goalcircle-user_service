package entity

import "time"

type User struct {
	Id           string  `gorm:"primaryKey"`
	FullName     string
	Email        string
	Password     string
	GoogleAuthId string
	UserType     string 
	CreatedAt    time.Time
	UpdatedAt    time.Time
}


