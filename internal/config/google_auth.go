package config

import (
	"errors"
	"fmt"
	"os"
	"time"
)

type GoogleAuthConfig struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
	TimeOut      time.Duration
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

	oauthSessionTimeOutStr := os.Getenv("GOOGLE_SESSION_TIMEOUT")

	oauthSessionTimeout := 10 * time.Minute // default
	if oauthSessionTimeOutStr != "" {
		var err error
		oauthSessionTimeout, err = time.ParseDuration(oauthSessionTimeOutStr)
		if err != nil {
			cb.errors = append(cb.errors, fmt.Errorf("GOOGLE_SESSION_TIMEOUT must be a valid duration (e.g. 5m, 10m): %w", err))
			return cb
		}
	}

	if len(cb.errors) > 0 {
		return cb
	}

	cb.config.GoogleAuthConfig = &GoogleAuthConfig{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		TimeOut:      oauthSessionTimeout,
	}

	return cb
}
