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
		Password:        request.Password,
		ConfirmPassword: request.ConfirmPassword,
	})

	if err != nil {
		return nil, domain.GRPCStatus(err)
	}

	return dt.ToRegisterResponse(response), nil
}

func (uh *authHandler) VerifyOtp(ctx context.Context, req *pb.VerifyOtpReq) (*pb.VerifyOtpRes, error) {
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
		Email: request.Email,
		Otp:   request.Otp,
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

	res, err := uh.authUseCase.Login(ctx, &ucdtos.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		return nil, domain.GRPCStatus(err)
	}

	return dt.ToLoginRes(res), nil
}

func (uh *authHandler) ResendOtp(ctx context.Context, pb *pb.ResendOtpReq) (*pb.ResendOtpRes, error) {

	ctx, cancel := context.WithTimeout(ctx, *uh.timeout)
	defer cancel()

	request := dt.ToResendOtpReq(pb)
	if validationErrs := uh.validater.Validation(request); validationErrs != nil {
		stWithDetails, err := ValidationError(validationErrs)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	res, err := uh.authUseCase.ResendOtp(ctx, &ucdtos.ResendOtpReq{
		Email:   request.Email,
		OtpType: request.OtpType,
	})

	if err != nil {
		return nil, domain.GRPCStatus(err)
	}
	return dt.ToResentOtpRes(res), nil
}

func (uh *authHandler) ForgotPassword(ctx context.Context, pb *pb.ForgotPasswordReq) (*pb.ForgotPasswordRes, error) {
	ctx, cancel := context.WithTimeout(ctx, *uh.timeout)
	defer cancel()

	request := dt.ToForgotPasswordReq(pb)

	if validationErrs := uh.validater.Validation(request); validationErrs != nil {
		stWithDetails, err := ValidationError(validationErrs)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	res, err := uh.authUseCase.ForgotPassword(ctx, &ucdtos.ForgotPasswordReq{
		Email: request.Email,
	})
	if err != nil {
		return nil, domain.GRPCStatus(err)
	}

	return dt.ToForgotPasswordRes(res), nil
}

func (uh *authHandler) VerifyForgotPassword(ctx context.Context, pb *pb.VerifyForgotPasswordReq) (*pb.VerifyForgotPasswordRes, error) {
	ctx, cancel := context.WithTimeout(ctx, *uh.timeout)
	defer cancel()

	request := dt.ToVerifyForgotPasswordOtpReq(pb)

	if validationErrs := uh.validater.Validation(request); validationErrs != nil {
		stWithDetails, err := ValidationError(validationErrs)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	res, err := uh.authUseCase.VerifyForgotPasswordOtp(ctx, &ucdtos.VerifyForgotPasswordOtpReq{
		Email: request.Email,
		Otp:   request.Otp,
	})

	if err != nil {
		return nil, domain.GRPCStatus(err)
	}

	return dt.ToVerifyForgotPasswordOtpRes(res), nil
}

func (uh *authHandler) ResetPassword(ctx context.Context, pb *pb.ResetPasswordReq) (*pb.ResetPasswordRes, error) {
	ctx, cancel := context.WithTimeout(ctx, *uh.timeout)
	defer cancel()

	request := dt.ToResetPasswordReq(pb)

	if validationErrs := uh.validater.Validation(request); validationErrs != nil {
		stWithDetails, err := ValidationError(validationErrs)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	res, err := uh.authUseCase.ResetPassword(ctx, &ucdtos.ResetPasswordReq{
		Email:    request.Email,
		Password: request.Password,
	})

	if err != nil {
		return nil, domain.GRPCStatus(err)
	}
	return dt.ToResetPasswordRes(res), nil
}

func (uh *authHandler) RenewAcccessToken(ctx context.Context, pb *pb.RenewAccessTokenReq) (*pb.RenewAccessTokenRes, error) {
	ctx, cancel := context.WithTimeout(ctx, *uh.timeout)
	defer cancel()

	request := dt.ToRenewAccessTokenReq(pb)

	if validationErrs := uh.validater.Validation(request); validationErrs != nil {
		stWithDetails, err := ValidationError(validationErrs)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	res, err := uh.authUseCase.RenewAccessToken(ctx, &ucdtos.RenewAcccessTokenReq{
		RefreshToken: request.RefreshToken,
	})
	if err != nil {
		return nil, domain.GRPCStatus(err)
	}
	return dt.ToRenewAccessTokenRes(res), nil
}

func (uh *authHandler) LogOut(ctx context.Context, pb *pb.LogOutReq) (*pb.LogOutRes, error) {
	ctx, cancel := context.WithTimeout(ctx, *uh.timeout)
	defer cancel()

	request := dt.ToLogOutReq(pb)

	if validationErrs := uh.validater.Validation(request); validationErrs != nil {
		stWithDetails, err := ValidationError(validationErrs)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}
	res, err := uh.authUseCase.LogOut(ctx, &ucdtos.LogOutReq{
		RefreshToken: request.RefreshToken,
	})
	if err != nil {
		return nil, domain.GRPCStatus(err)
	}
	return dt.ToLogoutRes(res), nil

}

func (uh *authHandler) OnboardingAddRole(ctx context.Context, pb *pb.OnboardingAddRoleReq) (*pb.OnboardingAddRoleRes, error) {
	ctx, cancel := context.WithTimeout(ctx, *uh.timeout)
	defer cancel()

	request := dt.ToOnboardingRoleReq(pb)

	if validationErrs := uh.validater.Validation(request); validationErrs != nil {
		stWithDetails, err := ValidationError(validationErrs)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	res, err := uh.authUseCase.OnboardingAddRole(ctx, &ucdtos.OnboardingRoleReq{
		UserId: request.UserId,
		Role:   request.UserRole,
	})

	if err != nil {
		return nil, domain.GRPCStatus(err)
	}

	return dt.ToOnboardingRoleRes(res), nil

}

func (h *authHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenReq) (*pb.ValidateTokenRes, error) {
	claims, err := h.authUseCase.ValidateToken(ctx, req.Token)
	if err != nil {
		return nil, domain.GRPCStatus(err)
	}
	return &pb.ValidateTokenRes{
		UserId:  claims.ID,
		Email:   claims.Email,
		Role:    claims.Role,
		IsValid: true,
	}, nil
}
