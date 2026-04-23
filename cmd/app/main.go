package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	cnfg "github.com/Junaidmdv/goalcircle-user_service/internal/config"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/repository"
	authHandler "github.com/Junaidmdv/goalcircle-user_service/internal/handler/grpc/auth"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/bycrypt"
	psql "github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/persistence/postgres"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/redis"
	sr "github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/server"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/otp"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/uid"
	at "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/auth"
	logger "github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/tokens"
	vl "github.com/Junaidmdv/goalcircle-user_service/pkg/validater"
	auth_pb "github.com/Junaidmdv/goalcircle-user_service/proto/pb"
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
		Build()
	logger.Info("configration is done")
	if errs != nil {
		for _, err := range errs {
			logger.Error("configration error", "error", err)
			fmt.Println()
		}
		return
	}

	validater, err := vl.NewValidater()
	if err != nil {
		logger.Error("validation package initilisation error", "error", err)
		return
	}

	//user authentication

	//postgres connection
	datbaseConnectin, err := psql.NewDatabase(config.Postgres)
	if err != nil {
		logger.Error("database initilisation error", "error", err)
		return
	}
	if err = datbaseConnectin.Migration(); err != nil {
		logger.Error("database migration error", "error", err)
		return
	}

	userRepo := repository.NewUserRepository(datbaseConnectin.DB, logger)
	uidGenerater := uid.NewUUIDGenerater()
	otpService := otp.NewSMSOtpService(config.Twilio)
	redisClient := redis.NewRedisClient(config.Redis)
	sessionStore := repository.NewSessionStorage(redisClient.Client)
	hashingCost := 14
	passwordHashing := bycrypt.NewBycriptHasher(hashingCost, logger)
	token := tokens.NewTokenMaker(config.JWT)
	authusecase := at.NewAuthUsecase(userRepo, logger, &config.GRPC.TimeOut, uidGenerater, otpService, passwordHashing, token, sessionStore)

	auth_handler := authHandler.NewAuthHandler(authusecase, logger, validater, &config.GRPC.TimeOut)
	server := sr.NewGrpcServer()
	auth_pb.RegisterAuthServiceServer(server.Server, auth_handler)

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
