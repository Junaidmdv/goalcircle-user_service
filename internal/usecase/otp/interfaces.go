package otp

import (
	"context"
	uc_dtos "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"
)

type OtpUsecase interface {
	ResendOtp(context.Context, *uc_dtos.ResendOtpReq) (*uc_dtos.ResendOtpResponse, error)
	VerifyOtp(context.Context, *uc_dtos.VerifyOtpRequest) (*uc_dtos.VerifyOtpResponse, error)
}
