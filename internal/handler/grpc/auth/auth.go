package auth

import (
	"context"
	"time"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain"
	dt "github.com/Junaidmdv/goalcircle-user_service/internal/handler/grpc/dtos"
	uc "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/auth"
	ucdtos "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
	vl "github.com/Junaidmdv/goalcircle-user_service/pkg/validater"
	"github.com/Junaidmdv/goalcircle-user_service/proto/pb"
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

func (uh *authHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
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
		PhoneNum:        req.PhoneNum,
		Password:        request.Password,
		ConfirmPassword: request.ConfirmPassword,
	})

	if err != nil {
		return nil, domain.GRPCStatus(err)
	}

	return dt.ToRegisterResponse(response), nil
}

func (uh *authHandler) VerifyOtp(ctx context.Context, req *pb.OtpReq) (*pb.OtpRes, error) {
	ctx, cancel := context.WithTimeout(ctx, *uh.timeout)
	defer cancel()

	request := dt.ToVerifyOtpReq(req)
	if validationErrs := uh.validater.Validation(request); validationErrs != nil {
		stWithDetails, err := ValidationError(validationErrs)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	res, err := uh.authUseCase.VerifyOtp(ctx, &ucdtos.VerifyOtpRequest{
		Email:    request.Email,
		PhoneNum: request.PhoneNum,
		Otp:      request.Otp,
	})

	if err != nil {
		return nil, domain.GRPCStatus(err)
	}

	return dt.ToVerifyOtpRes(res), nil
}

func (uh *authHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, *uh.timeout)
	defer cancel()

	request := dt.ToLoginRequest(req)
	if validationErrs := uh.validater.Validation(request); validationErrs != nil {
		stWithDetails, err := ValidationError(validationErrs)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	} 

	return nil, nil
}


func(uh *authHandler)ResendOtp(ctx context.Context,pb *pb.ResendOtpReq)(*pb.ResendOtpRes,error){
	return nil,nil
}


