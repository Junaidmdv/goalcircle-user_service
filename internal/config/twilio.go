package config

import (
	"errors"
	"os"
)

type TwilioConfig struct {
	AccountSID string
	AuthToken  string
}

func (cb *configBuilder) WithTwilio() ConfigBuilder {
	accsid := os.Getenv("TWILIO_ACCOUNTSID")
	if accsid == "" {
		cb.errors = append(cb.errors, errors.New("missing required env var: TWILIO_ACCOUNTSID"))
		return cb
	}

	authToken := os.Getenv("TWILIO_AUTHTOKEN")
	if authToken == "" {
		cb.errors = append(cb.errors, errors.New("missing required env var: TWILIO_AUTHTOKEN"))
		return cb
	}

	cb.config.Twilio.AccountSID = accsid
	cb.config.Twilio.AuthToken = authToken

	return cb
}
