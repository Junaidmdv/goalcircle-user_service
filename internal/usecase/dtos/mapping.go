package dtos

import "github.com/Junaidmdv/goalcircle-user_service/internal/domain/entity"

func ToRegisterResponse(res *entity.TempUser) *RegisterResponse {
	return &RegisterResponse{
		Email:     res.Email,
		PhoneNum:  res.PhoneNum,
		OtpStatus: true,
		OtpExpiry: res.ExpiresAt,
	}
}
