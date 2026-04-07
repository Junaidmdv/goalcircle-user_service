package auth

import (
	"context"
	"time"

	dt "github.com/junaidmdv/goalcirlcle/user_service/internal/handler/grpc/dtos"
	uc "github.com/junaidmdv/goalcirlcle/user_service/internal/usecase"
	ucdtos "github.com/junaidmdv/goalcirlcle/user_service/internal/usecase/dtos"
	"github.com/junaidmdv/goalcirlcle/user_service/pkg/logger"
	vl "github.com/junaidmdv/goalcirlcle/user_service/pkg/validater"
	"github.com/junaidmdv/goalcirlcle/user_service/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authHandler struct {
	pb.UnimplementedAuthServiceServer
	authUseCase uc.AuthUsecase
	logger      *logger.ZapLogger
	validater   *vl.Validater
}

func NewAuthHandler(aus uc.AuthUsecase, logger *logger.ZapLogger, validate *vl.Validater) *authHandler {
	return &authHandler{
		authUseCase: aus,
		logger:      logger,
		validater:   validate,
	}
}

func (uh *authHandler) RegisterUser(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	c, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	request := dt.ToRegisterReq(req)

	if validationErrs := uh.validater.Validation(request); validationErrs != nil {
		stWithDetails, err := ValidationError(validationErrs)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to attach details")
		}

		return nil, stWithDetails.Err()

	}

	response, err := uh.authUseCase.InitiateUserRegistration(c, &ucdtos.RegisterRequest{
		FullName:        request.FullName,
		Email:           request.Email,
		Password:        request.Password,
		ConfirmPassword: request.ConfirmPassword,
	})

	if err != nil {
		return nil, status.Error(codes.AlreadyExists, "already exist")
	}

	return dt.ToRegisterResponse(response), nil
}
