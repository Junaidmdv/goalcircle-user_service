package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/repository"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/bycrypt"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/twilio"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/uid"
	uc_dtos "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
	"go.uber.org/zap"
)

type authUsecase struct {
	userRepo     repository.UserRepository
	logger       logger.Logger
	timeout      *time.Duration
	uidGenerater uid.UuidGenerater
	otp          twilio.OtpService
	hash         bycrypt.PasswordHasher
}

func NewAuthUsecase(ur repository.UserRepository, logger logger.Logger, time *time.Duration, uidgen uid.UuidGenerater, otp twilio.OtpService, hash bycrypt.PasswordHasher) *authUsecase {
	return &authUsecase{
		userRepo:     ur,
		logger:       logger,
		timeout:      time,
		uidGenerater: uidgen,
		hash:         hash,
		otp:          otp,
	}
}

func (us *authUsecase) InitiateUserRegistration(ctx context.Context, input *uc_dtos.RegisterRequest) (*uc_dtos.RegisterResponse, error) {

	exist, err := us.userRepo.ExistByEmail(ctx, input.Email)
	if err != nil {
		us.logger.Error("internal error", zap.Error(err))
		return nil, domain.NewInternalError("something went wrong", err)
	}

	if exist {
		us.logger.Warn("dublicate email", errors.New("email already exist"))
		return nil, domain.NewConflictError("email already exist")
	}

	//phone number dublicate exist checking validation is added. But commented twilio is free tier only access verified number

	// exist, err = us.userRepo.ExistByPhoneNum(ctx, input.PhoneNum)
	// if err != nil {
	// 	us.logger.Error("internal error", zap.Error(err))
	// 	return nil, domain.NewInternalError("something went wrong", err)
	// }
	// if exist {
	// 	us.logger.Warn("dublicate email", errors.New("phone number already exist"))
	// 	return nil, domain.NewConflictError("phone number already exist")
	// }

	hashedPassword, err := us.hash.HashPassword(input.Password)
	if err != nil {
		us.logger.Error("failed to hash pasword", err)
		return nil, domain.NewInternalError("failed hash password", err)
	}

	otpres, err := us.otp.SendOtp(input.PhoneNum)
	if err != nil {
		return nil, us.otp.ParseTwilioError(err)
	}
	us.logger.Info("otp data", "data", otpres)

	res, err := us.userRepo.CreateTempUser(ctx, &entity.TempUser{ 
		FullName: input.FullName,
		Email:     input.Email,
		PhoneNum:  input.PhoneNum,
		Password:  hashedPassword,
		OTP:       otpres.Otp,
		ExpiresAt: otpres.ExpiresAt,
	})

	if err != nil {
		return nil, domain.NewInternalError("something went wrong", err)
	}
	return uc_dtos.ToRegisterResponse(res), nil

}

func (us *authUsecase) VerifyOtp(ctx context.Context, input *uc_dtos.OtpReq) (bool, error) {   

	
   

	return false, nil

}
