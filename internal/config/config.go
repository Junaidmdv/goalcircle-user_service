package config

import (
	"log"

	"github.com/joho/godotenv"
)

//created config using builder design pattern

type Config struct {
	GRPC     *GRPCConfig
	Postgres *PostgresConfig
	JWT      *JWTConfig
	Redis    *RedisConfig
	Twilio   *TwilioConfig
	SMTP     *SMTPConfig
	DiscStorage    *DiscStorageConfig
	GoogleAuthConfig *GoogleAuthConfig
}

type configBuilder struct {
	config *Config
	errors []error
}

type ConfigBuilder interface {
	WithGrpc() ConfigBuilder
	WithPostgres() ConfigBuilder
	WithJWT() ConfigBuilder
	WithRedis() ConfigBuilder
	WithTwilio() ConfigBuilder
	WithSMTP() ConfigBuilder
	WithDiscStorage() ConfigBuilder  
	WithGoogleAuth()ConfigBuilder
	Build() (*Config, []error)
}

func LoadConfig() *configBuilder {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found, reading from environment")
	}
	return &configBuilder{
		config: &Config{},
	}
}

func (cb *configBuilder) Build() (*Config, []error) {
	if len(cb.errors) > 0 {
		return nil, cb.errors
	}

	return cb.config, nil
}
