package dtos

import "github.com/junaidmdv/goalcircle/user_service/internal/domain/entity"

func ToRegisterResponse(res *entity.TempUser) *RegisterResponse {
	return &RegisterResponse{
		UserId:    res.ID,
		Email:     res.Email,
		PhoneNum:  res.PhoneNum,
		OtpStatus: true,
	}
}
