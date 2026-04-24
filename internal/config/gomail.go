package config

import (
	"errors"
	"log"
	"os"
	"strconv"
)

type SMTPConfig struct {
	FromEmail string
	Password  string //mail app password is used
	ServerURL string
	PORT      int
}

func (cb *configBuilder) WithSMTP() ConfigBuilder {
	sc := &SMTPConfig{
		FromEmail: os.Getenv("SMTP_FROM_EMAIL"),
		Password:  os.Getenv("SMTP_PASSWORD"),
		ServerURL: os.Getenv("SMTP_SERVER_URL"),
	}

	if sc.FromEmail == "" {
		cb.errors = append(cb.errors, errors.New("smtp from email is required"))
		return cb
	}

	if sc.Password == "" {
		cb.errors = append(cb.errors, errors.New("smtp password is required"))
		return cb
	}

	if sc.ServerURL == "" {
		log.Printf("SMTP_SERVER_URL not set, using default: smtp.gmail.com")
		sc.ServerURL = "smtp.gmail.com"
	}

	portStr := os.Getenv("SMTP_PORT")
	if portStr == "" {
		log.Printf("SMTP_PORT not set, using default: 587")
		sc.PORT = 587
	} else {
		port, err := strconv.Atoi(portStr) // just validate it's a number
		if err != nil {
			log.Printf("SMTP_PORT invalid, using default: 587")
			sc.PORT = 587
		} else {
			sc.PORT = port
		}
	}

	cb.config.SMTP = sc
	return cb
}
