package repository

import (
	"context"
	"fmt"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	ExistByEmail(context.Context, string) (bool, error)
	ExistByPhoneNum(context.Context, string) (bool, error)
	CreateTempUser(context.Context, *entity.TempUser) (*entity.TempUser, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
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
	// if err := ur.db.WithContext(ctx).Create(&tempUser).Error; err != nil {
	// 	return nil, err
	// }
	// return tempUser, nil

	if err := ur.db.WithContext(ctx).Raw(`INSERT INTO temp_users(full_name,email,phone_num,password,otp,expires_at,deleted_at)
	       VALUES(?,?,?,?,?,?,?)  
		   ON CONFLICT(email) 
		   DO UPDATE SET 
		   full_name=EXCLUDED.full_name,
		   phone_num=EXCLUDED.phone_num, 
		   password=EXCLUDED.password, 
           otp=EXCLUDED.otp, 
		   expires_at=EXCLUDED.expires_at,
           deleted_at=null
		   `, tempUser.FullName, tempUser.Email, tempUser.PhoneNum, tempUser.Password, tempUser.OTP, tempUser.ExpiresAt).Scan(&tempUser).Error; err != nil {
		return nil, err
	}
	return tempUser, nil

}

func (ur *userRepository) ExistByPhoneNum(ctx context.Context, phone_num string) (bool, error) {
	var count int64
	err := ur.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("phone_num = ?", phone_num).
		Count(&count).
		Error
	if err != nil {
		return false, fmt.Errorf("failed to check phone number existence: %w", err)
	}
	return count > 0, nil
}
