package config

import (
	"errors"
	"os"
)

type GoogleAuthConfig struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
}

func (cb *configBuilder) WithGoogleAuth() ConfigBuilder {
	clientId := os.Getenv("GOOGLE_CLIENT_ID")
	if clientId == "" {
		cb.errors = append(cb.errors, errors.New("missing google client id"))
	}

	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRETE")
	if clientSecret == "" {
		cb.errors = append(cb.errors, errors.New("missing google client secrete"))
	}

	redirectUrl := os.Getenv("GOOGLE_REDIRECT_URL")

	if redirectUrl == "" {
		cb.errors = append(cb.errors, errors.New("missing redirect url"))

	}

	if len(cb.errors) > 0 {
		return cb
	}

	cb.config.GoogleAuthConfig = &GoogleAuthConfig{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}

	return cb
}
