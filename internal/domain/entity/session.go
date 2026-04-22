package entity

import "time"

type Session struct {
	ID           string    `redis:"session_id"`
	UserEmail    string    `redis:"email"`
	RefreshToken string    `redis:"refresh_token"`
	IsRevoked    bool      `redis:"is_revoked"`
	CreatedAt    time.Time `redis:"created_at"`
	ExpiresAt    time.Time `redis:"expires_at"`
}
