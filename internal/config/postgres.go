package config

import "os"

type PostgresConfig struct {
	Host     string
	Port     string
	DB_Name     string
	User     string
	Password string
	SSLMode  string
}

func (cb *configBuilder) WithPostgress() ConfigBuilder {
	pc := &PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DB_Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	cb.config.Postgres = pc
	return cb
}
  


