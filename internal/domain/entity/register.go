package entity

import "time"

type PendingRegistration struct {
	ID        string // session key for redis
	Email     string
	Password  string // already hashed
	OTP       string
	ExpiresAt time.Time
}
