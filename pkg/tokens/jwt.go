package tokens

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtMaker struct {
	secreteKey string
}

const leeweetime = time.Second * 5

func (j *JwtMaker) GenerateToken(id int, email string, role string, duration time.Duration) (string, error) {

	claims, err := NewTokenClaims(id, email, role, duration)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secreteKey))
}

func (j *JwtMaker) VerifyToken(tokenstr string) (*UserClaims, error) {

	token, err := jwt.ParseWithClaims(tokenstr, &UserClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.secreteKey), nil
	}, jwt.WithLeeway(leeweetime))

	if err != nil {
		return nil, fmt.Errorf("error parsing token %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid user claims")
	}

	return claims, nil
}
