package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/junaidmdv/goalcirlcle/user_service/internal/domain"
	"github.com/junaidmdv/goalcirlcle/user_service/internal/domain/entity"
	"github.com/junaidmdv/goalcirlcle/user_service/internal/domain/repository"
	"github.com/junaidmdv/goalcirlcle/user_service/internal/infrastructure/twilio"
	"github.com/junaidmdv/goalcirlcle/user_service/internal/infrastructure/uid"
	uc_dtos "github.com/junaidmdv/goalcirlcle/user_service/internal/usecase/dtos"
	"github.com/junaidmdv/goalcirlcle/user_service/pkg/logger"
	"go.uber.org/zap"
)

type authUsecase struct {
	userRepo     repository.UserRepository
	logger       logger.Logger
	timeout      *time.Duration
	uidGenerater uid.UuidGenerater
	otp          twilio.OtpService
}

func NewAuthUsecase(ur repository.UserRepository, logger logger.Logger, time *time.Duration, uidgen uid.UuidGenerater, otp twilio.OtpService) *authUsecase {
	return &authUsecase{
		userRepo:     ur,
		logger:       logger,
		timeout:      time,
		uidGenerater: uidgen,
	}
}

func (us *authUsecase) InitiateUserRegistration(ctx context.Context, input *uc_dtos.RegisterRequest) (*uc_dtos.RegisterResponse, error) {

	context, cancel := context.WithTimeout(ctx, *us.timeout)
	defer cancel()

	exist, err := us.userRepo.ExistByEmail(context, input.Email)
	if err != nil {
		us.logger.Error("internal error", zap.Error(err))
		return nil, domain.NewInternalError("something went wrong", err)
	}

	if exist {
		us.logger.Warn("dublicate email", errors.New("email already exist"))
		return nil, domain.NewConflictError("email already exist")

	}

	otpres, err := us.otp.SendOtp(input.PhoneNum)
	if err != nil {
		return nil, us.otp.ParseTwilioError(err)
	}

	us.userRepo.CreateTempUser(ctx, &entity.TempUser{
		ID:        us.uidGenerater.Generate(),
		Email:     input.Email,
		PhoneNum:  input.PhoneNum,
		Password:  input.Password,
		OTP:       otpres.Otp,
		ExpiresAt: otpres.ExpiresAt,
	})

	return nil, nil

}

func (us *authUsecase) VerifyOtp() {

}
