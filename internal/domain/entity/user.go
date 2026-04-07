package entity

import "time"

type User struct {
	Id           int
	Fullname     string
	Email        string
	Password     string
	GoogleauthId string
	UserType     string //user: organiser,team manager
	CreatedAt    time.Time
}
