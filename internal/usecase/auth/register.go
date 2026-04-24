package usecase

import (
	"context"
	"time"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/repository"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/bycrypt"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/otp"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/uid"
	uc_dtos "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/tokens"
)

type AuthUsecase interface {
	InitiateUserRegistration(context.Context, *uc_dtos.RegisterRequest) (*uc_dtos.RegisterResponse, error)
	VerifyOtp(context.Context, *uc_dtos.VerifyOtpRequest) (*uc_dtos.VerifyOtpResponse, error)
	Login(context.Context, *uc_dtos.LoginRequest) (*uc_dtos.LoginResponse, error)
	ResendOtp(context.Context, *uc_dtos.ResendOtpReq) (*uc_dtos.ResendOtpResponse, error)
}

type authUsecase struct {
	userRepo     repository.UserRepository
	logger       logger.Logger
	timeout      *time.Duration
	uidGenerater uid.UuidGenerater
	hash         bycrypt.PasswordHasher
	token        *tokens.JwtMaker
	session      repository.SessionStorage
	email        *otp.EmailService
}

func NewAuthUsecase(ur repository.UserRepository, logger logger.Logger, time *time.Duration, uidgen uid.UuidGenerater, hash bycrypt.PasswordHasher, token *tokens.JwtMaker, session repository.SessionStorage, email *otp.EmailService) AuthUsecase {
	return &authUsecase{
		userRepo:     ur,
		logger:       logger,
		timeout:      time,
		uidGenerater: uidgen,
		hash:         hash,
		token:        token,
		session:      session,
	}
}

func (us *authUsecase) InitiateUserRegistration(ctx context.Context, input *uc_dtos.RegisterRequest) (*uc_dtos.RegisterResponse, error) {

	if err := us.userRepo.CheckEmailExist(ctx, input.Email); err != nil {
		return nil, err
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
		return nil, domain.NewInternalError("internal server error", err)
	}

	res, err := us.userRepo.CreateOrUpdateTempUser(ctx, &entity.TempUser{
		FullName: input.FullName,
		Email:    input.Email,
		PhoneNum: input.PhoneNum,
		Password: hashedPassword,
	})

	if err != nil {
		return nil, err
	}

	otpres, err := us.email.SendOTP(input.Email)
	if err != nil {
		return nil, us.email.MapMailError(err)
	}
	us.logger.Info("otp sended to the user", "email", input.Email, "otp", otpres.Otp)

	otpdata, err := us.userRepo.AddOtpData(ctx, &entity.Otp{
		TempUserID: res.ID,
		Otp:        otpres.Otp,
		Type:       string(entity.Register),
		ExpiresAt:  time.Now().Add(otpres.Expiry),
	})
	return uc_dtos.ToRegisterResponse(res, otpdata), nil
}

func (us *authUsecase) VerifyOtp(ctx context.Context, input *uc_dtos.VerifyOtpRequest) (*uc_dtos.VerifyOtpResponse, error) {

	tempUser, err := us.userRepo.GetTempUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	otpRecord, err := us.userRepo.GetLatestOtpRecord(ctx, tempUser.ID)
	if err != nil {
		return nil, err
	}

	if time.Now().After(otpRecord.ExpiresAt) {
		return nil, domain.NewUnAuthenticatedError("OTP expired")
	}

	if otpRecord.Otp != input.Otp {
		return nil, domain.NewUnAuthenticatedError("OTP has expired. Please request a new one")
	}

	if otpRecord.Attempts == 5 {
		return nil, domain.NewUnAuthenticatedError("OTP has reached max attempt.Resend otp")
	}

	us.logger.Info("otp verified", "email", tempUser.Email)

	user, err := us.userRepo.CreateUser(ctx, &entity.User{
		Id:       us.uidGenerater.Generate(),
		FullName: tempUser.FullName,
		Email:    tempUser.Email,
		Password: tempUser.Password,
	})

	if err != nil {
		return nil, err
	}

	us.logger.Info("user created", "id", user.Id, "email", user.Email)

	accessToken, accessClaims, err := us.token.GenerateToken(user.Id, user.Email, "user", us.token.AccessTokenExpiry)
	if err != nil {
		us.logger.Error("failed to generater token", "error", err)
		return nil, domain.NewInternalError("Something went wrong.Please try again later.", err)
	}

	refreshToken, refreshClaims, err := us.token.GenerateToken(user.Id, user.Email, "user", us.token.RefreshTokenExpiry)
	if err != nil {
		us.logger.Error("failed to generater token", "error", err)
		return nil, domain.NewInternalError("Something went wrong.Please try again later.", err)
	}

	if err := us.session.SaveSession(ctx, "session:"+refreshClaims.ID, &entity.Session{
		ID:           refreshClaims.ID,
		UserEmail:    refreshClaims.Email,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		CreatedAt:    time.Now(),
		ExpiresAt:    refreshClaims.ExpiresAt.Time,
	}); err != nil {
		return nil, err
	}

	return &uc_dtos.VerifyOtpResponse{
		UserId:            user.Id,
		AccessToken:       accessToken,
		AceessTokenExpiry: accessClaims.ExpiresAt.Time,
	}, nil

}

func (us *authUsecase) Login(ctx context.Context, input *uc_dtos.LoginRequest) (*uc_dtos.LoginResponse, error) {
	user, err := us.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if err := us.hash.ComparePassword(user.Password, input.Password); err != nil {
		return nil, err
	}

	accessToken, accessClaims, err := us.token.GenerateToken(user.Id, user.Email, "user", us.token.AccessTokenExpiry)
	if err != nil {
		us.logger.Error("failed to generater token", "error", err)
		return nil, domain.NewInternalError("Something went wrong.Please try again later.", err)
	}

	refreshToken, refreshClaims, err := us.token.GenerateToken(user.Id, user.Email, "user", us.token.RefreshTokenExpiry)
	if err != nil {
		us.logger.Error("failed to generater token", "error", err)
		return nil, domain.NewInternalError("Something went wrong.Please try again later.", err)
	}

	if err := us.session.SaveSession(ctx, "session:"+refreshClaims.ID, &entity.Session{
		ID:           refreshClaims.ID,
		UserEmail:    refreshClaims.Email,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		CreatedAt:    time.Now(),
		ExpiresAt:    refreshClaims.ExpiresAt.Time,
	}); err != nil {
		return nil, err
	}
	return &uc_dtos.LoginResponse{
		UserId:            user.Id,
		AccessToken:       accessToken,
		AccessTokenExpiry: accessClaims.ExpiresAt.Time,
	}, nil
}

func (us *authUsecase) ResendOtp(ctx context.Context, input *uc_dtos.ResendOtpReq) (*uc_dtos.ResendOtpResponse, error) {
	res, err := us.userRepo.GetTempUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	otpres, err := us.email.SendOTP(input.Email)
	if err != nil {
		return nil, us.email.MapMailError(err)
	}
	us.logger.Info("otp data", "data", otpres)

	otpdata, err := us.userRepo.AddOtpData(ctx, &entity.Otp{
		TempUserID: res.ID,
		Otp:        otpres.Otp,
		Type:       string(entity.ResendOtp),
		ExpiresAt:  time.Now().Add(otpres.Expiry),
	})

	return &uc_dtos.ResendOtpResponse{
		Success:   true,
		OtpExpiry: otpdata.ExpiresAt,
	}, nil
}

func (us *authUsecase) ForgotPassword(ctx context.Context) {

}

func (us *authUsecase) RevokeSession()

// func(us *authUsecase)
