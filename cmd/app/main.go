package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	cnfg "github.com/junaidmdv/goalcirlcle/user_service/internal/config"
	authHandler "github.com/junaidmdv/goalcirlcle/user_service/internal/handler/grpc/auth"
	psql "github.com/junaidmdv/goalcirlcle/user_service/internal/infrastructure/persistence/postgres"
	sr "github.com/junaidmdv/goalcirlcle/user_service/internal/infrastructure/server"
	at "github.com/junaidmdv/goalcirlcle/user_service/internal/usecase/auth"
	logger "github.com/junaidmdv/goalcirlcle/user_service/pkg/logger"
	vl "github.com/junaidmdv/goalcirlcle/user_service/pkg/validater"
	auth_pb "github.com/junaidmdv/goalcirlcle/user_service/proto/pb"
	"go.uber.org/zap"
)

func main() {

	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
		return
	}

	config, errs := cnfg.LoadConfig().
		WithGrpc().
		WithPostgres().
		//WithJWT().
		//WithRedis().
		Build()

	if errs != nil {

		for _, err := range errs {
			logger.Error("configration error", zap.Error(err))
			fmt.Println()
		}
		return
	}

	validater, err := vl.NewValidater()
	if err != nil {
		logger.Error("validation package initilisation error", zap.Error(err))
		return
	}

	//user authentication

	//postgres connection
	datbaseConnectin, err := psql.NewDatabase(config.Postgres)
	if err != nil {
		logger.Error("database initilisation error", zap.String("error", err.Error()))
	}
	userRepo := psql.NewUserRepository(datbaseConnectin)

	//redis connection
	authusecase := at.NewAuthUsecase(userRepo, logger.Logger)
	auth_handler := authHandler.NewAuthHandler(authusecase, logger, validater)
	server := sr.NewGrpcServer()
	auth_pb.RegisterAuthServiceServer(server.Server, auth_handler)

	go func() {
		logger.Info(fmt.Sprintf("server running on %d", config.GRPC.Port))
		if err := server.Run(config.GRPC.Port); err != nil {
			logger.Error("server error", zap.String("message", err.Error()), zap.Time("time", time.Now()))
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signalChan
	logger.Info("received signal, shutting down", zap.String("signal", sig.String()))

	signal.Stop(signalChan)
	server.Server.GracefulStop()
}
