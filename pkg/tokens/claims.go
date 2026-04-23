package tokens

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	UserId    string    `json:"user_id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func NewTokenClaims(id string, email string, role string, duration time.Duration) (*UserClaims, error) {
	tokenId, err := uuid.NewRandom()

	if err != nil {
		return nil, errors.New("failed generate token id")
	}

	return &UserClaims{
		UserId:    id,
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenId.String(),
			Subject:   email,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}
