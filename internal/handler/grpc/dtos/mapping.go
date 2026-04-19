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
		PhoneNum:        res.PhoneNum,
		Password:        res.Password,
		ConfirmPassword: res.ConfirmPassword,
	}
}

func ToRegisterResponse(res *dtos.RegisterResponse) *pb.RegisterResponse {
	return &pb.RegisterResponse{
		Email:        res.Email,
		PhoneNum:     res.PhoneNum,
		OtpStatus:    &res.OtpStatus,
		OtpExpiresAt: timestamppb.New(res.OtpExpiry),
	}
}

func ToOtpReq(res *pb.OtpReq) *VerifyOtpReq {
	return &VerifyOtpReq{
		Email:    res.Email,
		PhoneNum: res.PhoneNum,
		Otp:      res.Otp,
	}
}

func ToOtpRes(res *pb.OtpRes) *pb.OtpRes {
	return &pb.OtpRes{
		Verified: res.Verified,
	}
}
