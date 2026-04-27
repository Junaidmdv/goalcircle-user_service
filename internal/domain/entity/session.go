package entity

type Session struct {
	ID           string `redis:"session_id"`
	UserEmail    string `redis:"email"`
	RefreshToken string `redis:"refresh_token"`
	IsRevoked    string `redis:"is_revoked"`
	CreatedAt    string `redis:"created_at"`
	ExpiresAt    string `redis:"expires_at"`
}

