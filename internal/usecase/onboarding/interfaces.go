package onboarding

import (
	"context"

	uc_dtos "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"
)

type OnboardingUsecase interface {
	OnboardingAddRole(context.Context, *uc_dtos.OnboardingRoleReq) (*uc_dtos.OnboardingRoleRes, error)
	OnboardingAddTeamDetails(context.Context, *uc_dtos.OnboardingTeamDtlsReq) (*uc_dtos.OnboardingTeamDtlsRes, error)
	OnboardingAddOrganiserDetails(context.Context, *uc_dtos.OnboardingOrganiserDtlsReq) (*uc_dtos.OnboardingAddOrganiserDtlsRes, error)
}
