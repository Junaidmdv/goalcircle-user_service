package usermanagement

import (
	"context"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/repository"
	uc_dtos "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
)

type userManagementUsecase struct {
	repo   repository.AdminUserManagementRepository
	logger logger.Logger
}

func NewUserManagementUsecase(repo repository.AdminUserManagementRepository, logger logger.Logger) UserManagementUsecase {
	return &userManagementUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *userManagementUsecase) BlockUser(ctx context.Context, input *uc_dtos.BlockUserReq) (*uc_dtos.BlockUserRes, error) {
	IsBlock, err := uc.repo.IsBlockedUser(ctx, input.UserId)
	if err != nil {
		return nil, err
	}

	if IsBlock {
		return nil, domain.NewConflictError("user already blocked")
	}

	if err := uc.repo.BlockUser(ctx, input.UserId); err != nil {
		return nil, err
	}

	return &uc_dtos.BlockUserRes{
		Success: true,
	}, nil
}

func (uc *userManagementUsecase) UnBlockUser(ctx context.Context, input *uc_dtos.UnblockUserReq) (*uc_dtos.UnblockUserRes, error) {
	IsBlock, err := uc.repo.IsBlockedUser(ctx, input.UserId)
	if err != nil {
		return nil, err
	}

	if !IsBlock {
		return nil, domain.NewConflictError("user already unblocked")
	}

	if err := uc.repo.UnBlockUser(ctx, input.UserId); err != nil {
		return nil, err
	}

	return &uc_dtos.UnblockUserRes{
		Success: true,
	}, nil
}
func (uc *userManagementUsecase) GetUsers(ctx context.Context) ([]*uc_dtos.GetUserRes, error) {

	users, err := uc.repo.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	var userList []*uc_dtos.GetUserRes

	for _, user := range users {
		userList = append(userList, &uc_dtos.GetUserRes{
			UserId:   user.ID,
			Email:    user.Email, 
			IsBlocked: user.IsBlocked, 
			CreatedAt: user.CreatedAt,
		})
	}

	return userList, nil
}
