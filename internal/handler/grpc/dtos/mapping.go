package dtos

import (
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"
	"github.com/Junaidmdv/goalcircle-user_service/proto/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToRegisterReq(res *pb.RegisterRequest) *RegisterRequest {
	return &RegisterRequest{
		FullName:        res.FullName,
		Email:           res.Email,
		Password:        res.Password,
		ConfirmPassword: res.ConfirmPassword,
	}
}

func ToRegisterResponse(res *dtos.RegisterResponse) *pb.RegisterResponse {
	return &pb.RegisterResponse{
		Email:        res.Email,
		OtpStatus:    &res.OtpStatus,
		OtpExpiresAt: timestamppb.New(res.OtpExpiry),
	}
}

func ToVerifyOtpReq(res *pb.VerifyOtpReq) *VerifyOtpReq {
	return &VerifyOtpReq{
		Email: res.Email,
		Otp:   res.Otp,
	}
}

func ToVerifyOtpRes(res *dtos.VerifyOtpResponse) *pb.VerifyOtpRes {
	return &pb.VerifyOtpRes{
		Success:            true,
		SessionId:          res.SessionId,
		UserId:             res.UserId,
		FullName:           res.FullName,
		Email:              res.Email,
		AccessToken:        res.AccessToken,
		AccessTokenExpiry:  timestamppb.New(res.AceessTokenExpiry),
		RefreshToken:       res.RefreshToken,
		RefreshTokenExpiry: timestamppb.New(res.RefreshTokenExpiry),
	}
}

func ToLoginRequest(req *pb.LoginRequest) *LoginReq {
	return &LoginReq{
		Email:    req.Email,
		Password: req.Password,
	}
}

func ToLoginRes(res *dtos.LoginResponse) *pb.LoginResponse {
	return &pb.LoginResponse{
		UserId:             res.UserId,
		FullName:           res.FullName,
		Email:              res.Email,
		SussionId:          res.SessionId,
		AccessToken:        res.AccessToken,
		AccessTokenExpiry:  timestamppb.New(res.AccessTokenExpiry),
		RefreshToken:       res.RefreshToken,
		RefreshTokenExpiry: timestamppb.New(res.RefreshTokenExpiry),
	}
}

func ToResendOtpReq(req *pb.ResendOtpReq) *ResendOtpReq {
	return &ResendOtpReq{
		Email:   req.Email,
		OtpType: entity.OtpType(req.OtpType),
	}
}

func ToResentOtpRes(res *dtos.ResendOtpResponse) *pb.ResendOtpRes {
	return &pb.ResendOtpRes{
		Success: true,
	}
}

func ToForgotPasswordReq(res *pb.ForgotPasswordReq) *ForgotPasswordReq {
	return &ForgotPasswordReq{
		Email: res.Email,
	}
}

func ToForgotPasswordRes(res *dtos.ForgotPasswordRes) *pb.ForgotPasswordRes {
	return &pb.ForgotPasswordRes{
		Success: res.Success,
	}
}

func ToVerifyForgotPasswordOtpReq(pb *pb.VerifyForgotPasswordReq) *VerifyForgotPasswordOtpReq {
	return &VerifyForgotPasswordOtpReq{
		Email: pb.Email,
		Otp:   pb.Otp,
	}
}

func ToVerifyForgotPasswordOtpRes(res *dtos.VerifyForgotPasswordOtpRes) *pb.VerifyForgotPasswordRes {
	return &pb.VerifyForgotPasswordRes{
		Success:          res.Success,
		ResetToken:       res.ResetToken,
		ResetTokenExpiry: timestamppb.New(res.ExpiresAt),
	}
}

func ToResetPasswordReq(req *pb.ResetPasswordReq) *ResetPasswordReq {
	return &ResetPasswordReq{
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	}
}

func ToResetPasswordRes(res *dtos.ResetPasswordRes) *pb.ResetPasswordRes {
	return &pb.ResetPasswordRes{
		Success: res.Success,
	}
}

func ToRenewAccessTokenReq(req *pb.RenewAccessTokenReq) *RenewAcccessTokenReq {
	return &RenewAcccessTokenReq{
		RefreshToken: req.RefreshToken,
	}
}

func ToRenewAccessTokenRes(res *dtos.RenewAccessTokenRes) *pb.RenewAccessTokenRes {
	return &pb.RenewAccessTokenRes{
		AccessToken:       res.AccessToken,
		AccessTokenExpiry: res.AccessToken,
	}
}

func ToLogOutReq(req *pb.LogOutReq) *LogOutReq {
	return &LogOutReq{
		RefreshToken: req.RefreshToken,
	}
}

func ToLogoutRes(res *dtos.LogOutRes) *pb.LogOutRes {
	return &pb.LogOutRes{
		Success: res.Success,
	}
}

func ToOnboardingRoleReq(res *pb.OnboardingAddRoleReq) *OnboardRoleReq {
	return &OnboardRoleReq{
		UserId:   res.UserId,
		UserRole: res.Role,
	}
}

func ToOnboardingRoleRes(res *dtos.OnboardingRoleRes) *pb.OnboardingAddRoleRes {
	return &pb.OnboardingAddRoleRes{
		Success: res.Success,
	}
}
  


