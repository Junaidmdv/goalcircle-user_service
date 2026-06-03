package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
	"gorm.io/gorm"
)

type AdminUserManagementRepository interface {
	IsBlockedUser(context.Context, string) (bool, error)
	BlockUser(context.Context, string) error
	UnBlockUser(context.Context, string) error 
	GetUsers(context.Context)([]*entity.User,error)
}

type adminUserManagementRepository struct {
	db      *gorm.DB
	logger  logger.Logger
	timeout time.Duration
}

func NewAdminUserManagementRepository(db *gorm.DB, logger logger.Logger, time time.Duration) AdminUserManagementRepository {
	return &adminUserManagementRepository{
		db:      db,
		logger:  logger,
		timeout: time,
	}
}

func (am *adminUserManagementRepository) IsBlockedUser(ctx context.Context, userId string) (bool, error) {

	var user entity.User
	err := am.db.WithContext(ctx).Select("is_blocked").Where("id = ?", userId).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {

		return false, domain.NewNotFoundError("user not found")
	}
	if err != nil {
		am.logger.Error("database error", "error", err, "method", "adminUserManagementRepository.IsBlockUser")
		return false, domain.NewInternalError("Something went wrong please try again later.", err)
	}

	return user.IsBlocked, nil
}

func (am *adminUserManagementRepository) BlockUser(ctx context.Context, userId string) error {

	if err := am.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userId).Updates(map[string]interface{}{
		"is_blocked": true,
		"updated_at": time.Now(),
	}).Error; err != nil {
		am.logger.Error("database error", "error", err, "method", "adminUserManagementRepository.BlockUser")
		return domain.NewInternalError("Something went wrong.Please try again later.", err)
	}

	return nil
}

func (am *adminUserManagementRepository) UnBlockUser(ctx context.Context, userId string) error {

	if err := am.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userId).Updates(map[string]interface{}{
		"is_blocked": false,
		"updated_at": time.Now(),
	}).Error; err != nil {
		am.logger.Error("database error", "error", err, "method", "adminUserManagementRepository.BlockUser")
		return domain.NewInternalError("Something went wrong. Please try again later.", err)
	}

	return nil
}

func (am *adminUserManagementRepository) GetUsers(ctx context.Context) ([]*entity.User, error) {

	var users []*entity.User
	if err := am.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, domain.NewInternalError("Something went wrong. Please try again later", err)
	}
	return users, nil
}


