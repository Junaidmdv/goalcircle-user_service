package dtos

import (
	"github.com/junaidmdv/goalcirlcle/user_service/internal/usecase/dtos"
	"github.com/junaidmdv/goalcirlcle/user_service/proto/pb"
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
		SessionId: res.UserId,
		Email:     res.Email,
		OtpStatus: &res.OtpStatus,
	}
}
