package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func (cb *configBuilder) WithRedis() ConfigBuilder {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		cb.errors = append(cb.errors, errors.New("missing required env var: REDIS_HOST"))
	}

	port := os.Getenv("REDIS_PORT")
	if port == "" {
		cb.errors = append(cb.errors, errors.New("missing required env var: REDIS_PORT"))
	}

	db := 0
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		var err error
		db, err = strconv.Atoi(dbStr)
		if err != nil {
			cb.errors = append(cb.errors, fmt.Errorf("invalid REDIS_DB %q: must be an integer", dbStr))
		} else if db < 0 {
			cb.errors = append(cb.errors, fmt.Errorf("invalid REDIS_DB %d: must be non-negative", db))
		}
	}

	if host == "" || port == "" {
		return cb
	}

	cb.config.Redis = &RedisConfig{
		Host:     host,
		Port:     port,
		Password: os.Getenv("REDIS_PASSWORD"), 
		DB:       db,                        
	}

	return cb
}
