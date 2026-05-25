package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	cnfg "github.com/Junaidmdv/goalcircle-user_service/internal/config"
	sr "github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/server"
	logger "github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
)

func main() {

	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Fprintf(os.Stderr, "logger sync error: %v\n", err)
		}
	}()

	config, errs := cnfg.LoadConfig().
		WithGrpc().
		WithPostgres().
		WithTwilio().
		WithJWT().
		WithRedis().
		WithSMTP().
		WithGoogleAuth().
		WithDiscStorage().
		Build()
	logger.Info("configration is done")
	if errs != nil {
		for _, err := range errs {
			logger.Error("configration error", "error", err)
			fmt.Println()
		}
		return
	}

	server := sr.NewGrpcServer(logger, config)
	if err := server.BootStrapSetup(); err != nil {
		os.Exit(1)
	}
	go func() {
		logger.Info("server running", "port", config.GRPC.Port)
		if err := server.Run(config.GRPC.Port); err != nil {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signalChan
	logger.Info("received signal, shutting down", "signal", sig.String())

	signal.Stop(signalChan)
	server.Server.GracefulStop()
}
