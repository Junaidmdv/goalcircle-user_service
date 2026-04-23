package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
	"gorm.io/gorm"
)

type UserRepository interface {
	CheckEmailExist(context.Context, string) error
	ExistByPhoneNum(context.Context, string) error
	CreateOrUpdateTempUser(context.Context, *entity.TempUser) (*entity.TempUser, error)
	GetTempUserByEmail(context.Context, string) (*entity.TempUser, error)
	AddOtpData(context.Context, *entity.Otp) (*entity.Otp, error)
	GetLatestOtpRecord(context.Context, uint) (*entity.Otp, error)
	CreateUser(context.Context, *entity.User) (*entity.User, error)
	GetUserByEmail(context.Context, string) (*entity.User, error)
}

type userRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewUserRepository(db *gorm.DB, logger logger.Logger) UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

func (ur *userRepository) CheckEmailExist(ctx context.Context, email string) error {
	var exists bool

	err := ur.db.WithContext(ctx).
		Model(&entity.User{}).
		Select("1").
		Where("email = ?", email).
		Limit(1).
		Scan(&exists).
		Error

	if err != nil {
		ur.logger.Error("databse error", "error", err)
		return domain.NewInternalError("Something went wrong. Please try again later. ", err)
	}

	if exists {
		ur.logger.Error("dublicate account", "error", err)
		return domain.NewConflictError("email account already exist")
	}

	return nil
}

func (ur *userRepository) CreateOrUpdateTempUser(ctx context.Context, tempUser *entity.TempUser) (*entity.TempUser, error) {
	if err := ur.db.WithContext(ctx).Raw(`INSERT INTO temp_users(full_name,email,phone_num,password,deleted_at)
	       VALUES(?,?,?,?,NULL)  
		   ON CONFLICT(email) 
		   DO UPDATE SET 
		   full_name=EXCLUDED.full_name,
		   phone_num=EXCLUDED.phone_num, 
		   password=EXCLUDED.password, 
           deleted_at=null 
		   RETURNING *
		   `, tempUser.FullName, tempUser.Email, tempUser.PhoneNum, tempUser.Password).Scan(&tempUser).Error; err != nil {
		return nil, domain.NewInternalError("Something went wrong. Please try again later.", err)
	}
	return tempUser, nil

}

func (ur *userRepository) AddOtpData(ctx context.Context, otp *entity.Otp) (*entity.Otp, error) {

	if err := ur.db.WithContext(ctx).Create(otp).Error; err != nil {
		ur.logger.Error("database error", "error", err)
		return nil, domain.NewInternalError("Something went wrong. Please try again later", err)
	}
	return otp, nil

}

func (ur *userRepository) ExistByPhoneNum(ctx context.Context, phone_num string) error {
	var count int64
	err := ur.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("phone_num = ?", phone_num).
		Count(&count).
		Error
	if err != nil {
		ur.logger.Error("databse error", "error", err)
		return domain.NewInternalError("Something went wrong. Please try again later. ", err)
	}

	if count > 0 {
		ur.logger.Error("dublicate account", "error", err)
		return domain.NewConflictError("phone number is already exist")
	}

	return nil
}

func (ur *userRepository) GetTempUserByEmail(ctx context.Context, email string) (*entity.TempUser, error) {
	var user entity.TempUser

	if err := ur.db.WithContext(ctx).
		Model(&entity.TempUser{}).
		Where("email=?", email).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ur.logger.Info("record not found", "email", email, "error", err)
			return nil, domain.NewNotFoundError("No account register in this email.Please register first")
		}
		return nil, domain.NewInternalError("Something went wrong. Please try again later.", err)
	}
	return &user, nil
}

func (ur *userRepository) GetLatestOtpRecord(ctx context.Context, id uint) (*entity.Otp, error) {
	var otpRec entity.Otp
	if err := ur.db.WithContext(ctx).Where("email=?", id).Last(&otpRec).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ur.logger.Warn("record not found", "ID", id, "error", err)
			return nil, domain.NewNotFoundError("user otp record not found")
		}
		ur.logger.Error("database error", "err", err)
		return nil, domain.NewInternalError("database error", err)
	}
	return &otpRec, nil
}

func (ur *userRepository) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {

	if err := ur.db.WithContext(ctx).Create(user).Error; err != nil {
		ur.logger.Error("database error", "error", fmt.Errorf("falied to create user data:%v", err), "data", user)
		return nil, domain.NewInternalError("Something went wrong.Please try again later.", fmt.Errorf("failed to add user data:%v", err))
	}

	return user, nil
}

func (ur *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {

	var user entity.User
	if err := ur.db.WithContext(ctx).
		Where("email=?", email).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ur.logger.Warn("user email not found", "email", email)
			return nil, domain.NewNotFoundError("No account found with this email. Please register first.")
		}
		ur.logger.Error("database error", "error", err)
		return nil, domain.NewInternalError("Something went wrong. Please try again later.", err)

	}
	return &user, nil
}
