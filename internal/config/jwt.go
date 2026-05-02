package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

type JWTConfig struct {
	SecretKey       string
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration
	ResetTokenExp   time.Duration
}

func (cb *configBuilder) WithJWT() ConfigBuilder {
	jc := &JWTConfig{
		SecretKey: os.Getenv("JWT_SECRETKEY"),
	}
	if jc.SecretKey == "" {
		cb.errors = append(cb.errors, errors.New("jwt secret key is required"))
		return cb
	}

	accessTokenExpStr := os.Getenv("ACCESS_TOKEN_EXPIRYTIME")
	if accessTokenExpStr == "" {
		log.Printf("ACCESS_TOKEN_EXPIRYTIME not set, using default: 15m")
		accessTokenExpStr = "15m"
	}
	accessDuration, err := time.ParseDuration(accessTokenExpStr)
	if err != nil {
		cb.errors = append(cb.errors, fmt.Errorf("invalid ACCESS_TOKEN_EXPIRYTIME: %w", err))
		return cb
	}
	jc.AccessTokenExp = accessDuration

	refreshTokenExpStr := os.Getenv("REFRESH_TOKEN_EXPIRYTIME")
	if refreshTokenExpStr == "" {
		log.Printf("REFRESH_TOKEN_EXPIRYTIME not set, using default: 7d")
		refreshTokenExpStr = "168h"
	}
	refreshDuration, err := time.ParseDuration(refreshTokenExpStr)
	if err != nil {
		cb.errors = append(cb.errors, fmt.Errorf("invalid REFRESH_TOKEN_EXPIRYTIME: %w", err))
		return cb
	}
	jc.RefreshTokenExp = refreshDuration

	resetPasswordTokenExpStr := os.Getenv("RESET_TOKEN_EXPIRY")
	if refreshTokenExpStr == "" {
		log.Printf("REFRESH_TOKEN_EXPIRYTIME not set, using default: 5m")
		resetPasswordTokenExpStr = "5m"
	}
	resetPasswordTokenDuration, err := time.ParseDuration(resetPasswordTokenExpStr)
	if err != nil {
		cb.errors = append(cb.errors, fmt.Errorf("invalid REFRESH_TOKEN_EXPIRYTIME: %w", err))
		return cb
	}
	jc.ResetTokenExp = resetPasswordTokenDuration

	cb.config.JWT = jc
	return cb
}
