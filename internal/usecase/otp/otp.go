package otp

import (
	"context"
	"strconv"
	"time"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/repository"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/otp"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/uid"
	uc_dtos "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/tokens"
)

type otpUsecase struct {
	userRepo     repository.UserRepository
	uidGenerater uid.UuidGenerater
	logger       logger.Logger
	session      repository.SessionStorage
	token        *tokens.JwtMaker
	email        *otp.EmailService
}

func NewOtpUsecase(ur repository.UserRepository, uid uid.UuidGenerater, logger logger.Logger, session repository.SessionStorage, token *tokens.JwtMaker, email *otp.EmailService) OtpUsecase {
	return &otpUsecase{
		userRepo:     ur,
		uidGenerater: uid,
		logger:       logger,
		session:      session,
		token:        token,
		email:        email,
	}
}

func (us *otpUsecase) VerifyOtp(ctx context.Context, input *uc_dtos.VerifyOtpRequest) (*uc_dtos.VerifyOtpResponse, error) {

	tempUser, err := us.userRepo.GetTempUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	otpRecord, err := us.userRepo.GetLatestOtpRecord(ctx, tempUser.Email, entity.Register)
	if err != nil {
		return nil, err
	}

	us.logger.Info("otp data", "time", time.Now(), "data", otpRecord)

	if otpRecord.Attempts == entity.OtpMaxAttempts {
		return nil, domain.NewUnAuthenticatedError("OTP has reached max attempt.Resend otp")
	}

	if time.Now().After(otpRecord.ExpiresAt) {
		return nil, domain.NewUnAuthenticatedError("OTP expired")
	}

	if otpRecord.Otp != input.Otp {

		if err := us.userRepo.UpdateOtpAttempts(ctx, input.Email, entity.Register); err != nil {
			return nil, err
		}
		return nil, domain.NewUnAuthenticatedError("Invalid otp.Try again. ")
	}

	if err := us.userRepo.DeleteOtp(ctx, otpRecord.ID); err != nil {
		return nil, err
	}

	us.logger.Info("otp verified", "email", tempUser.Email)

	user, err := us.userRepo.CreateUser(ctx, &entity.User{
		ID:       us.uidGenerater.Generate(),
		Email:    tempUser.Email,
		Password: tempUser.Password,
		Role: entity.UNSPECIFIED,
	})

	if err != nil {
		return nil, err
	}

	accessToken, accessClaims, err := us.token.GenerateToken(user.ID, user.Email, entity.UNSPECIFIED, us.token.AccessTokenExpiry)
	if err != nil {
		us.logger.Error("failed to generater token", "error", err)
		return nil, domain.NewInternalError("Something went wrong.Please try again later.", err)
	}

	refreshToken, refreshClaims, err := us.token.GenerateToken(user.ID, user.Email, entity.UNSPECIFIED, us.token.RefreshTokenExpiry)
	if err != nil {
		us.logger.Error("failed to generater token", "error", err)
		return nil, domain.NewInternalError("Something went wrong.Please try again later.", err)
	}

	if err := us.session.SaveSession(ctx, "session:"+refreshClaims.ID, &entity.Session{
		ID:           refreshClaims.ID,
		Email:        refreshClaims.Email,
		Role:         refreshClaims.Role,
		RefreshToken: refreshToken,
		IsRevoked:    strconv.FormatBool(false),
		CreatedAt:    time.Now().Format(time.RFC3339),
		ExpiresAt:    refreshClaims.ExpiresAt.Time.Format(time.RFC3339),
	}); err != nil {
		return nil, err
	}

	return &uc_dtos.VerifyOtpResponse{
		SessionId:          refreshClaims.ID,
		UserId:             user.ID,
		Email:              user.Email,
		AccessToken:        accessToken,
		AceessTokenExpiry:  accessClaims.ExpiresAt.Time,
		RefreshToken:       refreshToken,
		RefreshTokenExpiry: refreshClaims.ExpiresAt.Time,
	}, nil

}

func (us *otpUsecase) ResendOtp(ctx context.Context, input *uc_dtos.ResendOtpReq) (*uc_dtos.ResendOtpResponse, error) {
	switch input.OtpType {
	case entity.Register:
		exist, err := us.userRepo.CheckEmailExistInTempUser(ctx, input.Email)
		if err != nil {
			return nil, err
		}
		if !exist {
			return nil, domain.NewNotFoundError("email not found. Please register first")
		}

	case entity.ForgotPassword:
		exist, err := us.userRepo.CheckEmailExist(ctx, input.Email)
		if err != nil {
			return nil, err
		}
		if !exist {
			return nil, domain.NewNotFoundError("email not found. Please use your registered email")
		}

	default:
		return nil, domain.NewBadRequestError("invalid OTP type")
	}

	otpRes, err := us.email.SendOTP(input.Email)
	if err != nil {
		return nil, us.email.MapMailError(err)
	}


	otpData, err := us.userRepo.AddOtpData(ctx, &entity.Otp{
		Email:     input.Email,
		Otp:       otpRes.Otp,
		Type:      string(input.OtpType),
		ExpiresAt: time.Now().Add(otpRes.Expiry),
	})
	if err != nil {
		return nil, err
	}

	return &uc_dtos.ResendOtpResponse{
		Success:   true,
		OtpExpiry: otpData.ExpiresAt,
	}, nil
}
