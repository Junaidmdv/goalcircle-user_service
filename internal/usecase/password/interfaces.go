package password

import (
	"context"
	uc_dtos "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"
)

type PasswordUsecase interface {
	VerifyForgotPasswordOtp(context.Context, *uc_dtos.VerifyForgotPasswordOtpReq) (*uc_dtos.VerifyForgotPasswordOtpRes, error)
	ForgotPassword(context.Context, *uc_dtos.ForgotPasswordReq) (*uc_dtos.ForgotPasswordRes, error)
	ResetPassword(context.Context, *uc_dtos.ResetPasswordReq) (*uc_dtos.ResetPasswordRes, error)
}
