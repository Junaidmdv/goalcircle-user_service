package config

import (
	"errors"
	"log"
	"os"
)

type DiscStorageConfig struct {
	BaseUrl  string
	BasePath string
}

func (cb *configBuilder) WithDiscStorage() ConfigBuilder {
	baseUrl := os.Getenv("DISC_STORAGE_BASE_URL")
	if baseUrl == "" {
		cb.errors = append(cb.errors, errors.New("DISC_STORAGE_BASE_URL is required"))
		return cb
	}

	basePath := os.Getenv("DISC_STORAGE_BASE_PATH")
	if basePath == "" {
		log.Print("DISC_STORAGE_BASE_PATH not set in env, using default: ./uploads")
		basePath = "./uploads"
	}

	cb.config.DiscStorage = &DiscStorageConfig{
		BaseUrl:  baseUrl,
		BasePath: basePath,
	}

	return cb
}
