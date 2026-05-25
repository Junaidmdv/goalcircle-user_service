package password

import (
	"context"
	"time"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/repository"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/bycrypt"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/otp"
	uc_dtos "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/tokens"
)

type passwordUsecase struct {
	userRepo repository.UserRepository
	email    *otp.EmailService
	logger   logger.Logger
	token    *tokens.JwtMaker
	hash     bycrypt.PasswordHasher
}

func NewPasswordUsecase(ur repository.UserRepository, email *otp.EmailService, logger logger.Logger, token *tokens.JwtMaker, hash bycrypt.PasswordHasher) PasswordUsecase {
	return &passwordUsecase{
		userRepo: ur,
		email:    email,
		logger:   logger,
		token:    token,
		hash:     hash,
	}
}

func (uc *passwordUsecase) ForgotPassword(ctx context.Context, input *uc_dtos.ForgotPasswordReq) (*uc_dtos.ForgotPasswordRes, error) {
	exist, err := uc.userRepo.CheckEmailExist(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, domain.NewNotFoundError("email not found.User registered email")
	}

	otp, err := uc.email.SendOTP(input.Email)
	if err != nil {
		return nil, uc.email.MapMailError(err)
	}

	otpExpiredAt := time.Now().Add(otp.Expiry)

	otpdata, err := uc.userRepo.AddOtpData(ctx, &entity.Otp{
		Email:     input.Email,
		Otp:       otp.Otp,
		Type:      string(entity.ForgotPassword),
		Attempts:  0,
		ExpiresAt: otpExpiredAt,
	})
	if err != nil {
		return nil, err
	}

	uc.logger.Info("otp data added in db", "data", otpdata)

	return &uc_dtos.ForgotPasswordRes{
		Success:   true,
		ExpiresAt: otpExpiredAt,
	}, nil
}

func (uc *passwordUsecase) VerifyForgotPasswordOtp(ctx context.Context, input *uc_dtos.VerifyForgotPasswordOtpReq) (*uc_dtos.VerifyForgotPasswordOtpRes, error) {
	exist, err := uc.userRepo.CheckEmailExist(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, domain.NewNotFoundError("email not found.User registered email")
	}

	otpRecord, err := uc.userRepo.GetLatestOtpRecord(ctx, input.Email, entity.ForgotPassword)
	if err != nil {
		return nil, err
	}
	if otpRecord.Attempts == entity.OtpMaxAttempts {
		return nil, domain.NewUnAuthenticatedError("OTP has reached max attempt.Resend otp")
	}

	if time.Now().After(otpRecord.ExpiresAt) {
		return nil, domain.NewUnAuthenticatedError("OTP expired.Resent otp")
	}

	if otpRecord.Otp != input.Otp {

		if err := uc.userRepo.UpdateOtpAttempts(ctx, input.Email, entity.Register); err != nil {
			return nil, err
		}
		return nil, domain.NewUnAuthenticatedError("Invalid otp.")
	}

	if err := uc.userRepo.DeleteOtp(ctx, otpRecord.ID); err != nil {
		return nil, err
	}

	uc.logger.Info("forgot password otp verified", "email", input.Email)

	user, err := uc.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	token, resetTokeClaims, err := uc.token.GenerateToken(user.ID, user.Email, "reset", uc.token.ResetPasswordTokenExpirty)
	if err != nil {
		uc.logger.Error("failed generate token", "method", "verify forget password", "layer", "usecase")
		return nil, domain.NewInternalError("Something went wrong. Please try again later.", err)
	}

	return &uc_dtos.VerifyForgotPasswordOtpRes{
		Success:    true,
		ResetToken: token,
		ExpiresAt:  resetTokeClaims.ExpiresAt.Time,
	}, nil

}

func (uc *passwordUsecase) ResetPassword(ctx context.Context, input *uc_dtos.ResetPasswordReq) (*uc_dtos.ResetPasswordRes, error) {

	claims, err := uc.token.VerifyToken(input.ResetToken)
	if err != nil {
		uc.logger.Warn("invalid reset token", "method", "ResetPassword", "error", err)
		return nil, err
	}

	if claims.Role != "reset" {
		uc.logger.Warn("invalid otp", "error", "invalid role", "method", "reset password")
		return nil, domain.NewUnAuthenticatedError("invalid otp")
	}

	hashedPassword, err := uc.hash.HashPassword(input.Password)
	if err != nil {
		uc.logger.Error("Failed hash password", "method", "ResetPassword", "error", err)
		return nil, domain.NewInternalError("Something went wrong. Please try again later", err)
	}

	if err := uc.userRepo.UpdatePassword(ctx, claims.Email, hashedPassword); err != nil {
		return nil, err
	}

	return &uc_dtos.ResetPasswordRes{
		Success: true,
	}, nil
}
