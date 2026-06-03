package onboarding

import (
	"context"

	uc_dtos "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"
)

type OnboardingUsecase interface {
	AddUserRole(context.Context, *uc_dtos.AddUserRoleReq) (*uc_dtos.AddUserRoleRes, error)
}
