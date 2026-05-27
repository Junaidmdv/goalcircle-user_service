package adminauth

import (
	"context"
	"strconv"
	"time"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/repository"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/bycrypt"
	"github.com/Junaidmdv/goalcircle-user_service/internal/infrastructure/uid"
	uc_dtos "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/tokens"
)

type adminAuth struct {
	adminRepo    repository.AdminRepository
	logger       logger.Logger
	hash         bycrypt.PasswordHasher
	uidGenerater uid.UuidGenerater
	token        *tokens.JwtMaker
	session      repository.SessionStorage
}

func NewAdminAuth(repo repository.AdminRepository, hash bycrypt.PasswordHasher, uid uid.UuidGenerater, token *tokens.JwtMaker) AdminAuth {
	return &adminAuth{
		adminRepo:    repo,
		hash:         hash,
		uidGenerater: uid,
		token:        token,
	}
}

func (ad *adminAuth) Register(ctx context.Context, input *uc_dtos.AdminAuthRegisterReq) (*uc_dtos.AdminAuthRegisterRes, error) {
	hashedPassword, err := ad.hash.HashPassword(input.Password)
	if err != nil {
		ad.logger.Error("failed to hash pasword", err)
		return nil, domain.NewInternalError("internal server error", err)
	}

	res, err := ad.adminRepo.Create(ctx, &entity.Admin{
		FullName: input.FullName,
		Email:    input.Email,
		Password: hashedPassword,
	})

	return &uc_dtos.AdminAuthRegisterRes{
		AdminId: res.ID,
		Email:   res.Email,
	}, nil
}

func (ad *adminAuth) Login(ctx context.Context, input *uc_dtos.AdminLoginReq) (*uc_dtos.AdminLoginRes, error) {
	res, err := ad.adminRepo.GetAdminByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if err := ad.hash.ComparePassword(res.Password, input.Password); err != nil {
		return nil, err
	}
	accessToken, accessClaims, err := ad.token.GenerateToken(res.ID, res.Email, string(entity.ADMIN), ad.token.AccessTokenExpiry)
	if err != nil {
		ad.logger.Error("failed to generater token", "error", err)
		return nil, domain.NewInternalError("Something went wrong.Please try again later.", err)
	}

	refreshToken, refreshClaims, err := ad.token.GenerateToken(res.ID, res.Email, string(entity.ADMIN), ad.token.RefreshTokenExpiry)
	if err != nil {
		ad.logger.Error("failed to generater token", "error", err)
		return nil, domain.NewInternalError("Something went wrong.Please try again later.", err)
	}

	key := "session:" + refreshClaims.ID
	if err := ad.session.SaveSession(ctx, key, &entity.Session{
		ID:           refreshClaims.ID,
		Email:        refreshClaims.Email,
		Role:         entity.ADMIN,
		RefreshToken: refreshToken,
		IsRevoked:    strconv.FormatBool(false),
		CreatedAt:    time.Now().Format(time.RFC3339),
		ExpiresAt:    refreshClaims.ExpiresAt.Time.Format(time.RFC3339),
	}); err != nil {
		return nil, err
	}

	return &uc_dtos.AdminLoginRes{
		SessionId:          refreshClaims.ID,
		UserId:             res.ID,
		FullName:           res.FullName,
		Email:              res.Email,
		AccessToken:        accessToken,
		AceessTokenExpiry:  accessClaims.ExpiresAt.Time,
		RefreshToken:       refreshToken,
		RefreshTokenExpiry: refreshClaims.ExpiresAt.Time,
	}, nil
}
