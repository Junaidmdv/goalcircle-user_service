package onboarding

import (
	"context"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/repository"
	uc_dtos "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
)

type onboardingUsecase struct {
	userRepo    repository.UserRepository
	fileStorage repository.FileStorage 
	logger     logger.Logger
}

func NewOnboardingUsecase(ur repository.UserRepository, fs repository.FileStorage,logger logger.Logger) OnboardingUsecase {
	return &onboardingUsecase{
		userRepo:    ur,
		fileStorage: fs, 
		logger:logger,
	}

}

func (uc *onboardingUsecase) OnboardingAddRole(ctx context.Context, input *uc_dtos.OnboardingRoleReq) (*uc_dtos.OnboardingRoleRes, error) {

	if err := uc.userRepo.UpdateUserType(ctx, input.UserId, input.Role); err != nil {
		return nil, err
	}

	return &uc_dtos.OnboardingRoleRes{
		Success: true,
	}, nil
}

func (uc *onboardingUsecase) OnboardingAddTeamDetails(ctx context.Context, input *uc_dtos.OnboardingTeamDtlsReq) (*uc_dtos.OnboardingTeamDtlsRes, error) {

	return nil, nil
}

func (uc *onboardingUsecase) OnboardingAddOrganiserDetails(ctx context.Context, input *uc_dtos.OnboardingOrganiserDtlsReq) (*uc_dtos.OnboardingAddOrganiserDtlsRes, error) {
	return nil, nil
}


