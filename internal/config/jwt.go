package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

type JWTConfig struct {
	// SecretKey         string
	PriviteKeyPath  string
	PublicKeyPath   string
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration
	ResetTokenExp   time.Duration
}

func (cb *configBuilder) WithJWT() ConfigBuilder {

	// secretekey := os.Getenv("JWT_SECRETKEY")
	// if secretekey == "" {
	// 	cb.errors = append(cb.errors, errors.New("jwt secret key is required"))
	// }

	jwtprivetekeypath := os.Getenv("JWT_PRIVATE_KEY_PATH")
	if jwtprivetekeypath == "" {
		cb.errors = append(cb.errors, errors.New("jwt secret key is required"))
	}

	jwtpublickeypath := os.Getenv("JWT_PUBLIC_KEY_PATH")
	if jwtpublickeypath == "" {
		cb.errors = append(cb.errors, errors.New("jwt secret key is required"))
	}

	accessTokenExpStr := os.Getenv("ACCESS_TOKEN_EXPIRYTIME")
	if accessTokenExpStr == "" {
		log.Printf("ACCESS_TOKEN_EXPIRYTIME not set, using default: 15m")
		accessTokenExpStr = "15m"
	}
	accessDuration, err := time.ParseDuration(accessTokenExpStr)
	if err != nil {
		cb.errors = append(cb.errors, fmt.Errorf("invalid ACCESS_TOKEN_EXPIRYTIME: %w", err))
	}
	accessTokenExp := accessDuration

	refreshTokenExpStr := os.Getenv("REFRESH_TOKEN_EXPIRYTIME")
	if refreshTokenExpStr == "" {
		log.Printf("REFRESH_TOKEN_EXPIRYTIME not set, using default: 7d")
		refreshTokenExpStr = "168h"
	}
	refreshDuration, err := time.ParseDuration(refreshTokenExpStr)
	if err != nil {
		cb.errors = append(cb.errors, fmt.Errorf("invalid REFRESH_TOKEN_EXPIRYTIME: %w", err))
	}
	refreshTokenExp := refreshDuration

	resetPasswordTokenExpStr := os.Getenv("RESET_TOKEN_EXPIRY")
	if refreshTokenExpStr == "" {
		log.Printf("REFRESH_TOKEN_EXPIRYTIME not set, using default: 5m")
		resetPasswordTokenExpStr = "5m"
	}
	resetPasswordTokenDuration, err := time.ParseDuration(resetPasswordTokenExpStr)
	if err != nil {
		cb.errors = append(cb.errors, fmt.Errorf("invalid REFRESH_TOKEN_EXPIRYTIME: %w", err))
	}
	resetTokenExp := resetPasswordTokenDuration

	if len(cb.errors) > 0 {
		return cb
	}

	cb.config.JWT = &JWTConfig{
		PriviteKeyPath:  jwtprivetekeypath,
		PublicKeyPath:   jwtpublickeypath,
		AccessTokenExp:  accessTokenExp,
		RefreshTokenExp: refreshTokenExp,
		ResetTokenExp:   resetTokenExp,
	}

	return cb
}
