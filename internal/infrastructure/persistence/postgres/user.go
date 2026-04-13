package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/junaidmdv/goalcircle/user_service/internal/domain/entity"
	"github.com/junaidmdv/goalcircle/user_service/internal/domain/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db      *gorm.DB
	timeout *time.Duration
}

func NewUserRepository(db *gorm.DB, timeout time.Duration) repository.UserRepository {
	return &userRepository{
		db:      db,
		timeout: &timeout,
	}
}

func (ur *userRepository) ExistByEmail(ctx context.Context, email string) (bool, error) {
	context, cancel := context.WithTimeout(ctx, *ur.timeout)
	defer cancel()

	var count int64
	err := ur.db.WithContext(context).
		Model(&entity.User{}).
		Where("email = ?", email).
		Count(&count).
		Error
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return count > 0, nil
}

func (ur *userRepository) CreateTempUser(ctx context.Context, tempUser *entity.TempUser) (*entity.TempUser, error) {
	context, cancel := context.WithTimeout(ctx, *ur.timeout)
	defer cancel()

	if err := ur.db.WithContext(context).Create(&tempUser).Error; err != nil {
		return nil, err
	}
	return tempUser, nil
}
