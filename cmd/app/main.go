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
	"github.com/junaidmdv/goalcirlcle/user_service/internal/infrastructure/persistence/postgres"
	sr "github.com/junaidmdv/goalcirlcle/user_service/internal/infrastructure/server"
	at "github.com/junaidmdv/goalcirlcle/user_service/internal/usecase/auth"
	logger "github.com/junaidmdv/goalcirlcle/user_service/pkg/logger"
	vl "github.com/junaidmdv/goalcirlcle/user_service/pkg/validater"
	auth_pb "github.com/junaidmdv/goalcirlcle/user_service/proto/pb"
	"go.uber.org/zap"
)

func main() {

	config, errs := cnfg.LoadConfig().
		WithGrpc().
		WithPostgress().
		//WithJWT().
		WithRedis().
		Build()

	if errs != nil {
		log.Fatal(errs)
		return
	}

	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
		return
	}

	validater, err := vl.NewValidater()

	if err != nil {
		logger.Error("validation package initilisation error", zap.String("error", err.Error()))
		return
	}

	//user authentication  


    //postgres connection 
	datbaseConnectin,err:=postgres.NewDatabase(config.Postgres) 
	if err != nil{
		 logger.Error("database initilisation error",zap.String("error",err.Error()))
	}
	userRepo:=postgres.NewUserRepository(datbaseConnectin) 
 

	


   //redis connection 
	authusecase := at.NewAuthUsecase(userRepo,)
	auth_handler := authHandler.NewAuthHandler(authusecase, logger, validater)
	server := sr.NewGrpcServer()
	auth_pb.RegisterAuthServiceServer(server.Server, auth_handler)


	go func() {
		if err := server.Run(config.GRPC.Port); err != nil {
			logger.Error("server error", zap.String("message", err.Error()), zap.Time("time", time.Now()))
		}
		logger.Info(fmt.Sprintf("server running on %s", config.GRPC.Port))
	}()



	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	signal.Stop(signalChan)

}


