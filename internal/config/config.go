package config

//created config using builder design pattern

type Config struct {
	GRPC     *GRPCConfig
	Postgres *PostgresConfig
	JWT      *JWTConfig 
	Redis    *RedisConfig
}

type configBuilder struct {
	config *Config
	errors []error
}

type ConfigBuilder interface {
	WithGrpc() ConfigBuilder
	WithPostgress() ConfigBuilder
	WithJWT() ConfigBuilder  
	WithRedis()ConfigBuilder
	Build()(*configBuilder,[]error)
}

func LoadConfig() *configBuilder {
	return &configBuilder{
		config: &Config{},
	}
}

func (cb *configBuilder) Build() (*configBuilder, []error) {
	if len(cb.errors) > 0 {
		return nil, cb.errors
	}

	return cb, nil
}
