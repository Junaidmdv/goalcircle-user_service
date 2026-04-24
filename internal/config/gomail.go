package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type SMTPConfig struct {
	FromEmail  string
	Password   string //mail app password is used
	ServerURL  string
	Port       int
	OtpExpiry  time.Duration
	MaxAttempt int
}

func (cb *configBuilder) WithSMTP() ConfigBuilder {
	var errs []error

	fromEmail := os.Getenv("SMTP_FROM_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	if fromEmail == "" {
		errs = append(errs, errors.New("SMTP_FROM_EMAIL is required"))
	}
	if password == "" {
		errs = append(errs, errors.New("SMTP_PASSWORD is required"))
	}

	// Collect all errors before returning
	if len(errs) > 0 {
		cb.errors = append(cb.errors, errs...)
		return cb
	}

	serverURL := os.Getenv("SMTP_SERVER_URL")
	if serverURL == "" {
		log.Printf("SMTP_SERVER_URL not set, using default: smtp.gmail.com")
		serverURL = "smtp.gmail.com"
	}

	port := 587
	if portStr := os.Getenv("SMTP_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err != nil {
			log.Printf("SMTP_PORT invalid (%q), using default: 587", portStr)
		} else if p < 1 || p > 65535 {
			log.Printf("SMTP_PORT out of range (%d), using default: 587", p)
		} else {
			port = p
		}
	} else {
		log.Printf("SMTP_PORT not set, using default: 587")
	}

	timeoutStr := os.Getenv("OTP_EXPIRY_TIME")
	if timeoutStr == "" {
		cb.errors = append(cb.errors, errors.New("OTP_EXPIRY_TIME is required"))
		return cb
	}

	otpExpiry, err := time.ParseDuration(timeoutStr)
	if err != nil {
		cb.errors = append(cb.errors, fmt.Errorf("OTP_EXPIRY_TIME must be a valid duration (e.g. 5m, 10m): %w", err))
		return cb
	}
	if otpExpiry < 1*time.Minute {
		cb.errors = append(cb.errors, errors.New("OTP_EXPIRY_TIME must be at least 1 minute"))
		return cb
	}

	maxAttempt := 5
	maxAttemptsStr := os.Getenv("OTP_MAX_ATTEMP")
	if maxAttemptsStr == "" {
		log.Printf("Maximum attempt ott is not added,Using default 5")
	} else {
		val, err := strconv.Atoi(maxAttemptsStr)
		if err != nil {
			log.Printf("Invalid maximum attemp value,Using default 5")
		} else {
			maxAttempt = val
		}
	}

	cb.config.SMTP = &SMTPConfig{
		FromEmail:  fromEmail,
		Password:   password,
		ServerURL:  serverURL,
		Port:       port,
		OtpExpiry:  otpExpiry,
		MaxAttempt: maxAttempt,
	}

	return cb
}
