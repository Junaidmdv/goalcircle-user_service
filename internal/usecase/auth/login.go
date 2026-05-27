package auth

import (
	"context"
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
	}
}

func (us *authUsecase) Login(ctx context.Context, input *uc_dtos.LoginRequest) (*uc_dtos.LoginResponse, error) {
	user, err := us.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if err := us.hash.ComparePassword(user.Password, input.Password); err != nil {
		return nil, err
	}

	accessToken, accessClaims, err := us.token.GenerateToken(user.ID, user.Email,user.UserType, us.token.AccessTokenExpiry)
	if err != nil {
		us.logger.Error("failed to generater token", "error", err)
		return nil, domain.NewInternalError("Something went wrong.Please try again later.", err)
	}

	refreshToken, refreshClaims, err := us.token.GenerateToken(user.ID, user.Email, user.UserType, us.token.RefreshTokenExpiry)
	if err != nil {
		us.logger.Error("failed to generater token", "error", err)
		return nil, domain.NewInternalError("Something went wrong.Please try again later.", err)
	}

	key := "session:" + refreshClaims.ID
	if err := us.session.SaveSession(ctx, key, &entity.Session{
		ID:           refreshClaims.ID,
		Email:        refreshClaims.Email,
		Role:         entity.UNSPECIFIED,
		RefreshToken: refreshToken,
		IsRevoked:    strconv.FormatBool(false),
		CreatedAt:    time.Now().Format(time.RFC3339),
		ExpiresAt:    refreshClaims.ExpiresAt.Time.Format(time.RFC3339),
	}); err != nil {
		return nil, err
	}
	return &uc_dtos.LoginResponse{
		SessionId:          refreshClaims.ID,
		UserId:             user.ID,
		FullName:           user.FullName,
		Email:              user.Email,
		AccessToken:        accessToken,
		AccessTokenExpiry:  accessClaims.ExpiresAt.Time,
		RefreshToken:       refreshToken,
		RefreshTokenExpiry: refreshClaims.ExpiresAt.Time,
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

	if sessiondata.Email != claims.Email {
		uc.logger.Warn("invalid token", "error", errors.New("token email mismatch"))
		return nil, domain.NewUnAuthenticatedError("Invalid refresh token. Please login again")
	}

	accessToken, accessClaims, err := uc.token.GenerateToken(sessiondata.ID, sessiondata.Email, claims.Role, uc.token.AccessTokenExpiry)
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
