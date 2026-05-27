package oauth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/repository"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/oauth"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/uid"
	uc_dtos "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/tokens"
)

type oauthUsecase struct {
	googleOauth  *oauth.GoogleOauth
	userRepo     repository.UserRepository
	logger       logger.Logger
	timeout      *time.Duration
	uidGenerater uid.UuidGenerater
	token        *tokens.JwtMaker
	session      repository.SessionStorage
}

func NewOauthUsecase(googleOauth *oauth.GoogleOauth, ur repository.UserRepository, logger logger.Logger, timeOut time.Duration, uid uid.UuidGenerater, token *tokens.JwtMaker, session repository.SessionStorage) OauthUsecase {
	return &oauthUsecase{
		googleOauth:  googleOauth,
		userRepo:     ur,
		logger:       logger,
		timeout:      &timeOut,
		uidGenerater: uid,
		token:        token,
		session:      session,
	}
}

func generateState() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (uc *oauthUsecase) GoogleOauth(ctx context.Context, input *uc_dtos.GoogleOauthReq) (*uc_dtos.GoogleOauthRes, error) {
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

func (uc *oauthUsecase) GoogleOauthCallback(ctx context.Context, input *uc_dtos.GoogleCallbackReq) (*uc_dtos.GoogleCallbackRes, error) {
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
		ID:           uc.uidGenerater.Generate(),
		FullName:     authdata.Name,
		Email:        authdata.Email,
		GoogleAuthId: authdata.ID,
		UserType:     entity.UNSPECIFIED,
	})

	if err != nil {
		return nil, err
	}

	uc.logger.Info("user created", "id", user.ID, "email", user.Email)

	accessToken, accessClaims, err := uc.token.GenerateToken(user.ID, user.Email, "user", uc.token.AccessTokenExpiry)
	if err != nil {
		uc.logger.Error("failed to generater token", "error", err)
		return nil, domain.NewInternalError("Something went wrong.Please try again later.", err)
	}

	refreshToken, refreshClaims, err := uc.token.GenerateToken(user.ID, user.Email, "user", uc.token.RefreshTokenExpiry)
	if err != nil {
		uc.logger.Error("failed to generater token", "error", err)
		return nil, domain.NewInternalError("Something went wrong.Please try again later.", err)
	}

	if err := uc.session.SaveSession(ctx, "session:"+refreshClaims.ID, &entity.Session{
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

	return &uc_dtos.GoogleCallbackRes{
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
