package server

import (
	"fmt"
	"net"

	cnfg "github.com/Junaidmdv/goalcircle-user_service/internal/config"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/repository"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/bycrypt"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/disk"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/oauth"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/otp"
	psql "github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/persistence/postgres"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/redis"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/uid"
	"github.com/Junaidmdv/goalcircle-user_service/internal/usecase/auth"

	authHandler "github.com/Junaidmdv/goalcircle-user_service/internal/handler/grpc/auth"
	oauthUc "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/oauth"

	otpUc "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/otp"
	"github.com/Junaidmdv/goalcircle-user_service/internal/usecase/register"

	"github.com/Junaidmdv/goalcircle-user_service/internal/usecase/onboarding"
	"github.com/Junaidmdv/goalcircle-user_service/internal/usecase/password"

	"github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/tokens"
	vl "github.com/Junaidmdv/goalcircle-user_service/pkg/validater"
	auth_pb "github.com/Junaidmdv/goalcircle-user_service/proto/pb"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	Server *grpc.Server
	logger logger.Logger
	config *cnfg.Config
}

func NewGrpcServer(logger logger.Logger, config *cnfg.Config) *GRPCServer {
	server := grpc.NewServer()

	return &GRPCServer{
		Server: server,
		logger: logger,
		config: config,
	}
}

func (s *GRPCServer) BootStrapSetup() error {

	validater, err := vl.NewValidater()
	if err != nil {
		s.logger.Error("validation package initilisation error", "error", err)
		return err
	}

	//user authentication

	//postgres connection
	datbaseConnectin, err := psql.NewDatabase(s.config.Postgres)
	if err != nil {
		s.logger.Error("database initilisation error", "error", err)
		return err
	}
	if err = datbaseConnectin.Migration(); err != nil {
		s.logger.Error("database migration error", "error", err)
		return err
	}

	userRepo := repository.NewUserRepository(datbaseConnectin.DB, s.logger, s.config.GRPC.TimeOut)
	uidGenerater := uid.NewUUIDGenerater()
	//otpService := otp.NewSMSOtpService(config.Twilio)
	redisClient := redis.NewRedisClient(s.config.Redis)
	sessionStore := repository.NewSessionStorage(redisClient.Client)
	hashingCost := 14
	passwordHashing := bycrypt.NewBycriptHasher(hashingCost, s.logger)
	token := tokens.NewTokenMaker(s.config.JWT, s.logger)
	emailService, err := otp.NewEmailService(s.config.SMTP)
	if err != nil {
		s.logger.Error("failed setup otp service email", "error", err)
		return err
	}
	diskStorage, err := disk.NewDiskStorage(s.config.DiscStorage, s.logger)
	if err != nil {
		s.logger.Error("failed setup disk storage", "error", err)
		return err
	}

	googleOauthSetup := oauth.NewGoogleOauth(s.config.GoogleAuthConfig)

	passwordUsecase := password.NewPasswordUsecase(userRepo, emailService, s.logger, token, passwordHashing)
	oauthUsecas := oauthUc.NewOauthUsecase(googleOauthSetup, userRepo, s.logger, s.config.GRPC.TimeOut, uidGenerater, token, sessionStore)
	authUsecase := auth.NewAuthUsecase(userRepo, s.logger, &s.config.GRPC.TimeOut, uidGenerater, passwordHashing, token, sessionStore, emailService, googleOauthSetup)
	registerUsecase := register.NewRegisterUsecase(userRepo, passwordHashing, s.logger, emailService)
	onboardingUsecase := onboarding.NewOnboardingUsecase(userRepo, diskStorage, s.logger)
	otpUsecas := otpUc.NewOtpUsecase(userRepo, uidGenerater, s.logger, sessionStore, token, emailService)

	auth_handler := authHandler.NewAuthHandler(authUsecase, oauthUsecas, onboardingUsecase, passwordUsecase, registerUsecase, otpUsecas, s.logger, validater, &s.config.GRPC.TimeOut)

	auth_pb.RegisterAuthServiceServer(s.Server, auth_handler)

	return nil
}

func (s *GRPCServer) Run(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	return s.Server.Serve(lis)
}
