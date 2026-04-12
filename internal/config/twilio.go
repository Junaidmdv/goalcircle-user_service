package config

import (
	"errors"
	"log"
	"os"
	"strconv"
)

type TwilioConfig struct {
	AccountSID    string
	AuthToken     string
	FromNum       string
	OtpExpiryTime int
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

	fromNum := os.Getenv("TWILIO_FROMNUM")
	if fromNum == "" {
		cb.errors = append(cb.errors, errors.New("missing required env var: FROM_NUM"))
		return cb
	}

	otpExpiry := os.Getenv("OTP_EXPIRY_MINUTES")
	if otpExpiry == "" {
		log.Printf("OTP_EXPIRY_MINUTES not set, defaulting to 5 minutes")
		otpExpiry = "5"
	}

	expiry, err := strconv.Atoi(otpExpiry)
	if err != nil {
		log.Printf("invalid OTP_EXPIRY_MINUTES value %q, defaulting to 5 minutes", otpExpiry)
		expiry = 5
	}

	cb.config.Twilio = &TwilioConfig{
		AccountSID:    accsid,
		AuthToken:     authToken,
		FromNum:       fromNum,
		OtpExpiryTime: expiry,
	}

	return cb
}
