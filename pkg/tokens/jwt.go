package tokens

import (
	"errors"
	"fmt"
	"time"

	"github.com/Junaidmdv/goalcircle-user_service/internal/config"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
)

func NewTokenMaker(jwtcnfg *config.JWTConfig, logger logger.Logger) *JwtMaker {
	return &JwtMaker{
		secreteKey:                jwtcnfg.SecretKey,
		AccessTokenExpiry:         jwtcnfg.AccessTokenExp,
		RefreshTokenExpiry:        jwtcnfg.RefreshTokenExp,
		ResetPasswordTokenExpirty: jwtcnfg.ResetTokenExp,
		logger:                    logger,
	}
}

type JwtMaker struct {
	secreteKey                string
	AccessTokenExpiry         time.Duration
	RefreshTokenExpiry        time.Duration
	ResetPasswordTokenExpirty time.Duration
	logger                    logger.Logger
}

const leeweetime = time.Second * 5

func (j *JwtMaker) GenerateToken(id string, email string, role string, duration time.Duration) (string, *UserClaims, error) {

	claims, err := NewTokenClaims(id, email, role, duration)
	if err != nil {
		j.logger.Error("token error", "error", err)
		return "", nil, domain.NewInternalError("Something went wrong. Please try again later", err)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenstr, err := token.SignedString([]byte([]byte(j.secreteKey)))
	if err != nil {
		j.logger.Error("failed to generate token", "error", err)
		return "", nil, domain.NewInternalError("Something went wrong Please try again later.", err)
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
		j.logger.Warn("token verification failed", "error", err)
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domain.NewUnAuthenticatedError("token expired")
		}

		return nil, domain.NewUnAuthenticatedError("invalid token")
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		j.logger.Error("token verification failed", "error", fmt.Errorf("invalid user claims"))
		return nil, domain.NewUnAuthenticatedError("invalid token")
	}

	return claims, nil
}
