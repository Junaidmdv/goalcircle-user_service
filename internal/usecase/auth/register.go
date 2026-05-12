package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/repository"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/bycrypt"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/oauth"
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
	VerifyForgotPasswordOtp(context.Context, *uc_dtos.VerifyForgotPasswordOtpReq) (*uc_dtos.VerifyForgotPasswordOtpRes, error)
	ForgotPassword(context.Context, *uc_dtos.ForgotPasswordReq) (*uc_dtos.ForgotPasswordRes, error)
	ResetPassword(context.Context, *uc_dtos.ResetPasswordReq) (*uc_dtos.ResetPasswordRes, error)
	RenewAccessToken(context.Context, *uc_dtos.RenewAcccessTokenReq) (*uc_dtos.RenewAccessTokenRes, error)
	LogOut(context.Context, *uc_dtos.LogOutReq) (*uc_dtos.LogOutRes, error)
	OnboardingAddRole(context.Context, *uc_dtos.OnboardingRoleReq) (*uc_dtos.OnboardingRoleRes, error)
	OnboardingAddTeamDetails(context.Context, *uc_dtos.OnboardingTeamDtlsReq) (*uc_dtos.OnboardingTeamDtlsRes, error)
	OnboardingAddOrganiserDetails(context.Context, *uc_dtos.OnboardingOrganiserDtlsReq) (*uc_dtos.OnboardingAddOrganiserDtlsRes, error)
	ValidateToken(context.Context, string) (*tokens.UserClaims, error)
	GoogleOauth(context.Context, *uc_dtos.GoogleOauthReq) (*uc_dtos.GoogleOauthRes, error)
	GoogleOauthCallback(context.Context, *uc_dtos.GoogleCallbackReq) (*uc_dtos.GoogleCallbackRes, error)
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
	googleOauth  *oauth.GoogleOauth
}

func NewAuthUsecase(ur repository.UserRepository, logger logger.Logger, time *time.Duration, uidgen uid.UuidGenerater, hash bycrypt.PasswordHasher, token *tokens.JwtMaker, session repository.SessionStorage, email *otp.EmailService, googleOauth *oauth.GoogleOauth) AuthUsecase {
	return &authUsecase{
		userRepo:     ur,
		logger:       logger,
		timeout:      time,
		uidGenerater: uidgen,
		hash:         hash,
		token:        token,
		session:      session,
		email:        email,
		googleOauth:  googleOauth,
	}
}

func (us *authUsecase) InitiateUserRegistration(ctx context.Context, input *uc_dtos.RegisterRequest) (*uc_dtos.RegisterResponse, error) {

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

func (us *authUsecase) VerifyOtp(ctx context.Context, input *uc_dtos.VerifyOtpRequest) (*uc_dtos.VerifyOtpResponse, error) {

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
		Id:       us.uidGenerater.Generate(),
		FullName: tempUser.FullName,
		Email:    tempUser.Email,
		Password: tempUser.Password,
		UserType: entity.UNSPECIFIED,
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
		IsRevoked:    strconv.FormatBool(false),
		CreatedAt:    time.Now().Format(time.RFC3339),
		ExpiresAt:    refreshClaims.ExpiresAt.Time.Format(time.RFC3339),
	}); err != nil {
		return nil, err
	}

	return &uc_dtos.VerifyOtpResponse{
		SessionId:          refreshClaims.ID,
		UserId:             user.Id,
		FullName:           user.FullName,
		Email:              user.Email,
		AccessToken:        accessToken,
		AceessTokenExpiry:  accessClaims.ExpiresAt.Time,
		RefreshToken:       refreshToken,
		RefreshTokenExpiry: refreshClaims.ExpiresAt.Time,
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

	key := "session:" + refreshClaims.ID
	if err := us.session.SaveSession(ctx, key, &entity.Session{
		ID:           refreshClaims.ID,
		UserEmail:    refreshClaims.Email,
		RefreshToken: refreshToken,
		IsRevoked:    strconv.FormatBool(false),
		CreatedAt:    time.Now().Format(time.RFC3339),
		ExpiresAt:    refreshClaims.ExpiresAt.Time.Format(time.RFC3339),
	}); err != nil {
		return nil, err
	}
	return &uc_dtos.LoginResponse{
		SessionId:          refreshClaims.ID,
		UserId:             user.Id,
		FullName:           user.FullName,
		Email:              user.Email,
		AccessToken:        accessToken,
		AccessTokenExpiry:  accessClaims.ExpiresAt.Time,
		RefreshToken:       refreshToken,
		RefreshTokenExpiry: refreshClaims.ExpiresAt.Time,
	}, nil
}

func (us *authUsecase) ResendOtp(ctx context.Context, input *uc_dtos.ResendOtpReq) (*uc_dtos.ResendOtpResponse, error) {
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

	us.logger.Info("otp resent", "email", input.Email, "type", input.OtpType)

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

func (uc *authUsecase) ForgotPassword(ctx context.Context, input *uc_dtos.ForgotPasswordReq) (*uc_dtos.ForgotPasswordRes, error) {
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

func (uc *authUsecase) VerifyForgotPasswordOtp(ctx context.Context, input *uc_dtos.VerifyForgotPasswordOtpReq) (*uc_dtos.VerifyForgotPasswordOtpRes, error) {
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

	token, resetTokeClaims, err := uc.token.GenerateToken(user.Id, user.Email, "reset", uc.token.ResetPasswordTokenExpirty)
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

func (uc *authUsecase) ResetPassword(ctx context.Context, input *uc_dtos.ResetPasswordReq) (*uc_dtos.ResetPasswordRes, error) {

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

func (uc *authUsecase) RenewAccessToken(ctx context.Context, input *uc_dtos.RenewAcccessTokenReq) (*uc_dtos.RenewAccessTokenRes, error) {
	claims, err := uc.token.VerifyToken(input.RefreshToken)
	if err != nil {
		return nil, err
	}

	key := "session:" + claims.ID

	sessiondata, err := uc.session.GetSession(ctx, key)
	if err != nil {
		return nil, err
	}

	isRevoked, _ := strconv.ParseBool(sessiondata.IsRevoked)
	if isRevoked {
		uc.logger.Info("invalid token", "error", "user token revoked")
		return nil, domain.NewUnAuthenticatedError("Invalid refresh token.Please login again")
	}

	if sessiondata.UserEmail != claims.Email {
		uc.logger.Warn("invalid token", "error", errors.New("token email mismatch"))
		return nil, domain.NewUnAuthenticatedError("Invalid refresh token. Please login again")
	}

	accessToken, accessClaims, err := uc.token.GenerateToken(sessiondata.ID, sessiondata.UserEmail, claims.Role, uc.token.AccessTokenExpiry)
	if err != nil {
		return nil, err
	}

	return &uc_dtos.RenewAccessTokenRes{
		AccessToken:       accessToken,
		AccessTokenExpiry: accessClaims.ExpiresAt.Time,
	}, nil

}

func (uc *authUsecase) LogOut(ctx context.Context, input *uc_dtos.LogOutReq) (*uc_dtos.LogOutRes, error) {

	claims, err := uc.token.VerifyToken(input.RefreshToken)
	if err != nil {
		return nil, err
	}

	key := "session:" + claims.ID

	if err := uc.session.DeleteSession(ctx, key); err != nil {
		return nil, err
	}

	return &uc_dtos.LogOutRes{
		Success: true,
	}, nil

}

func (uc *authUsecase) ValidateToken(ctx context.Context, tokenStr string) (*tokens.UserClaims, error) {

	claims, err := uc.token.VerifyToken(tokenStr)
	if err != nil {
		return nil, err
	}

	session, err := uc.session.GetSession(ctx, "session:"+claims.ID)
	if err != nil {
		return nil, domain.NewUnAuthenticatedError("session not found")
	}

	isRevoked, _ := strconv.ParseBool(session.IsRevoked)

	if isRevoked {
		return nil, domain.NewUnAuthenticatedError("session has been revoked")
	}

	return claims, nil
}

func (uc *authUsecase) OnboardingAddRole(ctx context.Context, input *uc_dtos.OnboardingRoleReq) (*uc_dtos.OnboardingRoleRes, error) {

	if err := uc.userRepo.UpdateUserType(ctx, input.UserId, input.Role); err != nil {
		return nil, err
	}

	return &uc_dtos.OnboardingRoleRes{
		Success: true,
	}, nil
}

func (uc *authUsecase) OnboardingAddTeamDetails(ctx context.Context, input *uc_dtos.OnboardingTeamDtlsReq) (*uc_dtos.OnboardingTeamDtlsRes, error) {

	return nil, nil
}

func (uc *authUsecase) OnboardingAddOrganiserDetails(ctx context.Context, input *uc_dtos.OnboardingOrganiserDtlsReq) (*uc_dtos.OnboardingAddOrganiserDtlsRes, error) {
	return nil, nil
}

func generateState() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (uc *authUsecase) GoogleOauth(ctx context.Context, input *uc_dtos.GoogleOauthReq) (*uc_dtos.GoogleOauthRes, error) {
	goauth := uc.googleOauth.Config
	state, err := generateState()
	if err != nil {
		uc.logger.Error("failed generate state", "error", err)
		return nil, domain.NewInternalError("Something went wrong. Please try again later", err)
	}

	url := goauth.AuthCodeURL(state)

	return &uc_dtos.GoogleOauthRes{
		RedirectUrl: url,
		State:       state,
		ExpireAt:    time.Now().Add(uc.googleOauth.TimeOut),
	}, nil

}

func (uc *authUsecase) GoogleOauthCallback(ctx context.Context, input *uc_dtos.GoogleCallbackReq) (*uc_dtos.GoogleCallbackRes, error) {
	token, err := uc.googleOauth.Exchange(ctx, input.Code)
	if err != nil {
		return nil, err
	}

	authdata, err := uc.googleOauth.GetUserData(ctx, token)
	if err != nil {
		return nil, err
	}

	exist, err := uc.userRepo.CheckEmailExist(ctx, authdata.Email)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, domain.NewConflictError("email already exist")
	}

	user, err := uc.userRepo.CreateUser(ctx, &entity.User{
		Id:           uc.uidGenerater.Generate(),
		FullName:     authdata.Name,
		Email:        authdata.Email,
		GoogleAuthId: authdata.ID,
		UserType:     entity.UNSPECIFIED,
	})

	if err != nil {
		return nil, err
	}

	uc.logger.Info("user created", "id", user.Id, "email", user.Email)

	accessToken, accessClaims, err := uc.token.GenerateToken(user.Id, user.Email, "user", uc.token.AccessTokenExpiry)
	if err != nil {
		uc.logger.Error("failed to generater token", "error", err)
		return nil, domain.NewInternalError("Something went wrong.Please try again later.", err)
	}

	refreshToken, refreshClaims, err := uc.token.GenerateToken(user.Id, user.Email, "user", uc.token.RefreshTokenExpiry)
	if err != nil {
		uc.logger.Error("failed to generater token", "error", err)
		return nil, domain.NewInternalError("Something went wrong.Please try again later.", err)
	}

	if err := uc.session.SaveSession(ctx, "session:"+refreshClaims.ID, &entity.Session{
		ID:           refreshClaims.ID,
		UserEmail:    refreshClaims.Email,
		RefreshToken: refreshToken,
		IsRevoked:    strconv.FormatBool(false),
		CreatedAt:    time.Now().Format(time.RFC3339),
		ExpiresAt:    refreshClaims.ExpiresAt.Time.Format(time.RFC3339),
	}); err != nil {
		return nil, err
	}

	return &uc_dtos.GoogleCallbackRes{
		SessionId:          refreshClaims.ID,
		UserId:             user.Id,
		FullName:           user.FullName,
		Email:              user.Email,
		AccessToken:        accessToken,
		AccessTokenExpiry:  accessClaims.ExpiresAt.Time,
		RefreshToken:       refreshToken,
		RefreshTokenExpiry: refreshClaims.ExpiresAt.Time,
	}, nil
}
