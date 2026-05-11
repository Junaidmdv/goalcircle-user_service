package config

import (
	"errors"
	"os"
)

type GoogleAuthConfig struct {
	ClientId     string
	ClientSecret string
}

func (cb *configBuilder) WithGoogleAuth() ConfigBuilder {
	clientId := os.Getenv("GOOGLE_CLIENT_ID")
	if clientId == "" {
		cb.errors = append(cb.errors, errors.New("missing google client id"))
	}

	clientSecrete := os.Getenv("GOOGLE_CLIENT_SECRETE")
	if clientSecrete == "" {
		cb.errors = append(cb.errors, errors.New("missing google client secrete"))
	}

	if len(cb.errors) > 0 {
		return cb
	}

	googlecnfig := &GoogleAuthConfig{
		ClientId:     clientId,
		ClientSecret: clientSecrete,
	}

	cb.config.GoogleAuthConfig = googlecnfig

	return cb
}
