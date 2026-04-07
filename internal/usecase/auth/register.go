package usecase

import (
	"context"

	"github.com/junaidmdv/goalcirlcle/user_service/internal/domain/repository"
	uc_dtos"github.com/junaidmdv/goalcirlcle/user_service/internal/usecase/dtos"
	"go.uber.org/zap"
)

type authUsecase struct {
	userRepo   repository.UserRepository
	pendingUserRepo  repository.PendingUserRepository
	logger *zap.Logger
}

func NewAuthUsecase(ur repository.PendingUserRepository,pr repository.PendingUserRepository,logger *zap.Logger) *authUsecase {
	return &authUsecase{}
}

func (us *authUsecase) InitiateUserRegistration(ctx context.Context, input *uc_dtos.RegisterRequest) (*uc_dtos.RegisterResponse, error) { 
    
	return nil, nil

}

func (us *authUsecase) VerifyOtp() {

}
