package config

import "os"

type RedisConfig struct {
	Host string
	Port string
	DB   string
}

func (cb *configBuilder) WithRedis() ConfigBuilder {
	rc := &RedisConfig{
		Host: os.Getenv("REDIS_HOST"),
		Port: os.Getenv("REDIS_PORT"),
		DB:   os.Getenv("REDIS_DB"),
	}

	cb.config.Redis = rc

	return cb
}
