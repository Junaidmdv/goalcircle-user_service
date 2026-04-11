package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type GRPCConfig struct {
	Port    int
	TimeOut time.Duration
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

	timeoutStr := os.Getenv("TIMEOUT")
	timeout := time.Second * 5

	// timeoutStr := os.Getenv("TIMEOUT")
	// timeout := 5 * time.Second // default timeout

	if timeoutStr == "" {
		log.Print("TIMEOUT not set in env, using default: 5 seconds")
	} else {
		t, err := strconv.Atoi(timeoutStr)
		if err != nil {
			log.Printf("invalid TIMEOUT format %q: must be an integer", timeoutStr)
		} else if t < 2 || t > 10 {
			log.Printf("invalid TIMEOUT value %d: must be between 2 and 10", t)
		} else {
			timeout = time.Duration(t) * time.Second
		}
	}

	cb.config.GRPC = &GRPCConfig{
		Port:    val,
		TimeOut: timeout,
	}
	return cb
}
