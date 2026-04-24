package dtos

import "github.com/Junaidmdv/goalcircle-user_service/internal/domain/entity"

func ToRegisterResponse(res *entity.TempUser, otpdata *entity.Otp) *RegisterResponse {
	return &RegisterResponse{
		Email:     res.Email,
		OtpStatus: otpdata != nil,
		OtpExpiry: otpdata.ExpiresAt,
	}
}



