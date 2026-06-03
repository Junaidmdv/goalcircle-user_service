package auth

import (
	"context"

	uc_dtos "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/tokens"
)

type AuthUsecase interface {
	Login(context.Context, *uc_dtos.LoginRequest) (*uc_dtos.LoginResponse, error)
	RenewAccessToken(context.Context, *uc_dtos.RenewAcccessTokenReq) (*uc_dtos.RenewAccessTokenRes, error)
	LogOut(context.Context, *uc_dtos.LogOutReq) (*uc_dtos.LogOutRes, error)
	ValidateToken(context.Context, string) (*tokens.UserClaims, error) 
	ChangePassword(context.Context,*uc_dtos.ChangePasswordReq)(*uc_dtos.ChangePasswordRes,error)
}





