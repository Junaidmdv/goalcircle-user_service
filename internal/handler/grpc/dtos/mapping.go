package dtos

import (
	"github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"
	"github.com/Junaidmdv/goalcircle-user_service/proto/pb"
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
		UserId:    res.UserId,
		Email:     res.Email,
		OtpStatus: &res.OtpStatus,
	}
}
