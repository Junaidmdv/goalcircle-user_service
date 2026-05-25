package register

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
)

type registerUsecase struct {
	userRepo repository.UserRepository
	hash     bycrypt.PasswordHasher
	logger   logger.Logger
	email    *otp.EmailService
}

func NewRegisterUsecase(ur repository.UserRepository, hash bycrypt.PasswordHasher, logger logger.Logger, email *otp.EmailService) RegistrationUsecase {
	return &registerUsecase{
		userRepo: ur,
		hash:     hash,
		logger:   logger,
		email:    email,
	}
}

func (us *registerUsecase) InitiateUserRegistration(ctx context.Context, input *uc_dtos.RegisterRequest) (*uc_dtos.RegisterResponse, error) {

	exist, err := us.userRepo.CheckEmailExist(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, domain.NewConflictError("an account with this email already exists. Please sign in or use a different email")
	}

	hashedPassword, err := us.hash.HashPassword(input.Password)
	if err != nil {
		us.logger.Error("failed to hash pasword", err)
		return nil, domain.NewInternalError("internal server error", err)
	}

	res, err := us.userRepo.CreateOrUpdateTempUser(ctx, &entity.TempUser{
		FullName: input.FullName,
		Email:    input.Email,
		Password: hashedPassword,
	})

	if err != nil {
		return nil, err
	}

	otpres, err := us.email.SendOTP(input.Email)
	if err != nil {
		us.logger.Warn("use case error", "error", err)
		return nil, us.email.MapMailError(err)
	}
	us.logger.Info("otp sended to the user", "email", input.Email, "otp", otpres.Otp)

	otpdata, err := us.userRepo.AddOtpData(ctx, &entity.Otp{
		Email:     input.Email,
		Otp:       otpres.Otp,
		Type:      string(entity.Register),
		ExpiresAt: time.Now().Add(otpres.Expiry),
	})

	if err != nil {
		return nil, err
	}

	return &uc_dtos.RegisterResponse{
		Email:     res.Email,
		OtpStatus: otpdata != nil,
		OtpExpiry: otpdata.ExpiresAt,
	}, nil
}
