package adminauth

import (
	"context"
	"time"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain"
	"github.com/Junaidmdv/goalcircle-user_service/internal/handler/grpc/auth"
	dt "github.com/Junaidmdv/goalcircle-user_service/internal/handler/grpc/dtos"
	"github.com/Junaidmdv/goalcircle-user_service/internal/usecase/adminauth"
	"github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/validater"
	"github.com/Junaidmdv/goalcircle-user_service/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type adminAuthHandler struct {
	pb.UnimplementedAdminAuthServiceServer
	adminAuthUsecase adminauth.AdminAuthUsecase
	timeout          *time.Duration
	logger           logger.Logger
	validater        *validater.Validater
}

func NewAdminAuthHandler(aAuth adminauth.AdminAuthUsecase, timeout *time.Duration, logger logger.Logger, validater *validater.Validater) *adminAuthHandler {
	return &adminAuthHandler{
		adminAuthUsecase: aAuth,
		timeout:          timeout,
		logger:           logger,
		validater:        validater,
	}
}

func (au adminAuthHandler) AdminRegister(ctx context.Context, input *pb.AdminRegisterRequest) (*pb.AdminRegisterResponse, error) {
	context, cancel := context.WithTimeout(ctx, *au.timeout)
	defer cancel()

	request := dt.ToAdminRegisterReq(input)

	if validationErrs := au.validater.Validation(request); validationErrs != nil {
		stWithDetails, err := auth.ValidationError(validationErrs)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	res, err := au.adminAuthUsecase.Register(context, &dtos.AdminAuthRegisterReq{
		FullName: request.FullName,
		Email:    request.Email,
		Password: request.Password,
	})

	if err != nil {
		return nil, domain.GRPCStatus(err)
	}

	return &pb.AdminRegisterResponse{
		Email: res.Email,
	}, nil
}

func (au adminAuthHandler) AdminLogin(ctx context.Context, input *pb.AdminLoginRequest) (*pb.AdminLoginResponse, error) {
	context, cancel := context.WithTimeout(ctx, *au.timeout)
	defer cancel()

	request := dt.ToAdminLoginReq(input)

	if validationErrs := au.validater.Validation(request); validationErrs != nil {
		stWithDetails, err := auth.ValidationError(validationErrs)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	res, err := au.adminAuthUsecase.Login(context, &dtos.AdminLoginReq{
		Email:    request.Email,
		Password: request.Password,
	})

	if err != nil {
		return nil, domain.GRPCStatus(err)
	}

	return &pb.AdminLoginResponse{
		SessionId:          res.SessionId,
		AdminId:            res.AdminId,
		FullName:           res.FullName,
		AccessToken:        res.AccessToken,
		AccessTokenExpiry:  timestamppb.New(res.AceessTokenExpiry),
		RefreshToken:       res.RefreshToken,
		RefreshTokenExpiry: timestamppb.New(res.RefreshTokenExpiry),
	}, nil

}
