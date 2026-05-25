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
}




type OnboardingUsecase interface {
	OnboardingAddRole(context.Context, *uc_dtos.OnboardingRoleReq) (*uc_dtos.OnboardingRoleRes, error)
	OnboardingAddTeamDetails(context.Context, *uc_dtos.OnboardingTeamDtlsReq) (*uc_dtos.OnboardingTeamDtlsRes, error)
	OnboardingAddOrganiserDetails(context.Context, *uc_dtos.OnboardingOrganiserDtlsReq) (*uc_dtos.OnboardingAddOrganiserDtlsRes, error)
}



