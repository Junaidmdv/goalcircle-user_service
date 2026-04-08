package config

import (
	"errors"
	"os"
	"time"
)

type JWTConfig struct {
	SecretKey  string
	AccessExp  time.Duration
	RefreshExp time.Duration
}

func (cb *configBuilder) WithJWT() ConfigBuilder {
	jc := &JWTConfig{
		SecretKey: os.Getenv("JWT_SECRETEKEY"),
	}

	if jc.SecretKey == "" {
		cb.errors = append(cb.errors, errors.New("jwt secrete key is required"))
		return cb
	}
	cb.config.JWT = jc

	return cb

}
