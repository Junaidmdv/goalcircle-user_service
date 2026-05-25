package entity

import "time"

type User struct {
	ID           string `gorm:"primaryKey"`
	FullName     string
	Email        string
	Password     string
	GoogleAuthId string
	UserType     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Organiser    Organiser   `gorm:"foreignKey:UserID"`
	TeamManager  TeamManager `gorm:"foreignKey:UserID"`
}
