package postgres

import (
	"context"
	"fmt"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

// func (ur *userRepository) ExistByEmail(ctx context.Context, email string) (bool, error) {
// 	var count int64
// 	err := ur.db.WithContext(ctx).
// 		Model(&entity.User{}).
// 		Where("email = ?", email).
// 		Count(&count).
// 		Error
// 	if err != nil {
// 		return false, fmt.Errorf("failed to check email existence: %w", err)
// 	}
// 	return count > 0, nil
// }

func (ur *userRepository) ExistByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool

	err := ur.db.WithContext(ctx).
		Model(&entity.User{}).
		Select("1").
		Where("email = ?", email).
		Limit(1).
		Scan(&exists).
		Error

	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}

func (ur *userRepository) CreateTempUser(ctx context.Context, tempUser *entity.TempUser) (*entity.TempUser, error) {
	if err := ur.db.WithContext(ctx).Create(&tempUser).Error; err != nil {
		return nil, err
	}
	return tempUser, nil
}
