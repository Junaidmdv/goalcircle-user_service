package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type GRPCConfig struct {
	Port    int
	TimeOut string
}

func (cb *configBuilder) WithGrpc() ConfigBuilder {

	portStr := os.Getenv("GRPC_PORT")
	if portStr == "" {
		cb.errors = append(cb.errors, errors.New("grpc port is required"))
		return cb
	}

	val, err := strconv.Atoi(portStr)
	if err != nil {
		cb.errors = append(cb.errors, fmt.Errorf("invalid GRPC_PORT %q: must be an integer", portStr))
		return cb
	} else if val < 1 || val > 65535 {
		cb.errors = append(cb.errors, fmt.Errorf("invalid GRPC_PORT %d: must be between 1 and 65535", val))
		return cb

	}

	cb.config.GRPC = &GRPCConfig{
		Port: val,
	}
	return cb
}
