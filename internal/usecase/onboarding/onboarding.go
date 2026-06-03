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

func (uc *onboardingUsecase) AddUserRole(ctx context.Context, input *uc_dtos.AddUserRoleReq) (*uc_dtos.AddUserRoleRes, error) {

	if err := uc.userRepo.UpdateUserRole(ctx, input.UserId, input.Role); err != nil {
		return nil, err
	}

	return &uc_dtos.AddUserRoleRes{
		Success: true,
	}, nil
}

