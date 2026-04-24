package dtos

import (
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

func ToVerifyOtpReq(res *pb.OtpReq) *VerifyOtpReq {
	return &VerifyOtpReq{
		Email:    res.Email,
		Otp:      res.Otp,
	}
}

func ToVerifyOtpRes(res *dtos.VerifyOtpResponse) *pb.OtpRes {
	return &pb.OtpRes{
		UserId:            res.UserId,
		AccessToken:       res.AccessToken,
		AccessTokenExpiry: timestamppb.New(res.AceessTokenExpiry),
	}
}

func ToLoginRequest(req *pb.LoginRequest) *LoginRequest {
	return &LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}
}
