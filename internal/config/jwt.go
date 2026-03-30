package config

import (
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
	cb.config.JWT = jc

	return cb

}
 

