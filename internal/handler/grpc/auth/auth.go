package auth

import (
	"context"
	"time"

	dt "github.com/junaidmdv/goalcircle/user_service/internal/handler/grpc/dtos"
	uc "github.com/junaidmdv/goalcircle/user_service/internal/usecase"
	ucdtos "github.com/junaidmdv/goalcircle/user_service/internal/usecase/dtos"
	"github.com/junaidmdv/goalcircle/user_service/pkg/logger"
	vl "github.com/junaidmdv/goalcircle/user_service/pkg/validater"
	"github.com/junaidmdv/goalcircle/user_service/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authHandler struct {
	pb.UnimplementedAuthServiceServer
	authUseCase uc.AuthUsecase
	logger      logger.Logger
	validater   *vl.Validater
	timeout     *time.Duration
}

func NewAuthHandler(aus uc.AuthUsecase, logger logger.Logger, validate *vl.Validater, time *time.Duration) *authHandler {
	return &authHandler{
		authUseCase: aus,
		logger:      logger,
		validater:   validate,
		timeout:     time,
	}
}

func (uh *authHandler) RegisterUser(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	context, cancel := context.WithTimeout(ctx, *uh.timeout)
	defer cancel()

	request := dt.ToRegisterReq(req)

	if validationErrs := uh.validater.Validation(request); validationErrs != nil {
		stWithDetails, err := ValidationError(validationErrs)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	response, err := uh.authUseCase.InitiateUserRegistration(context, &ucdtos.RegisterRequest{
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
