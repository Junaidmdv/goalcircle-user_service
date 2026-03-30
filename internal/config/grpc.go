package config

import (
	"os"
)


type GRPCConfig struct {
	Port string
}

func (cb *configBuilder) WithGrpc() ConfigBuilder {

	gc := &GRPCConfig{
		Port: os.Getenv("GRPC_PORT"),
	}
	cb.config.GRPC = gc
	return cb
}
