package tokens

import (
	"fmt"
	"time"

	"github.com/Junaidmdv/goalcircle-user_service/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

func NewTokenMaker(jwtcnfg *config.JWTConfig) *JwtMaker {
	return &JwtMaker{
		secreteKey:         jwtcnfg.SecretKey,
		AccessTokenExpiry:  jwtcnfg.AccessTokenExp,
		RefreshTokenExpiry: jwtcnfg.RefreshTokenExp,
	}
}

type JwtMaker struct {
	secreteKey         string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

const leeweetime = time.Second * 5

func (j *JwtMaker) GenerateToken(id string, email string, role string, duration time.Duration) (string, *UserClaims, error) {

	claims, err := NewTokenClaims(id, email, role, duration)
	if err != nil {
		return "", nil, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenstr, err := token.SignedString([]byte([]byte(j.secreteKey)))
	if err != nil {
		return "", nil, fmt.Errorf("failed generate token %v", err)
	}
	return tokenstr, claims, nil
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
