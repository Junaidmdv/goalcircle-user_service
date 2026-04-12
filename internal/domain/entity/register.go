package entity

import "time"

type TempUser struct {
	ID        string // session key
	Email     string
	PhoneNum  string
	Password  string // already hashed
	OTP       string
	ExpiresAt time.Time
}
