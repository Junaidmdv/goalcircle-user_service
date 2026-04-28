package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
	"gorm.io/gorm"
)

type UserRepository interface {
	CheckEmailExist(context.Context, string) (bool, error)
	CreateOrUpdateTempUser(context.Context, *entity.TempUser) (*entity.TempUser, error)
	GetTempUserByEmail(context.Context, string) (*entity.TempUser, error)
	AddOtpData(context.Context, *entity.Otp) (*entity.Otp, error)
	GetLatestOtpRecord(context.Context, string, entity.OtpType) (*entity.Otp, error)
	CreateUser(context.Context, *entity.User) (*entity.User, error)
	GetUserByEmail(context.Context, string) (*entity.User, error)
	CheckEmailExistInTempUser(context.Context, string) (bool, error)
	UpdateOtpAttempts(context.Context, string, entity.OtpType) error
	UpdatePassword(context.Context, string, string) error
	DeleteOtp(context.Context, uint) error
	UpdateUserType(context.Context, string,string) error
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

func (ur *userRepository) CheckEmailExist(ctx context.Context, email string) (bool, error) {
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
		return false, domain.NewInternalError("Something went wrong. Please try again later. ", err)
	}

	if exists {
		return true, nil
	}

	return false, nil
}

func (ur *userRepository) CreateOrUpdateTempUser(ctx context.Context, tempUser *entity.TempUser) (*entity.TempUser, error) {
	if err := ur.db.WithContext(ctx).Raw(`INSERT INTO temp_users(full_name,email,phone_num,password,deleted_at)
	       VALUES(?,?,?,?,NULL)  
		   ON CONFLICT(email) WHERE deleted_at IS NULL
		   DO UPDATE SET 
		   full_name=EXCLUDED.full_name,
		   password=EXCLUDED.password, 
		   RETURNING *
		   `, tempUser.FullName, tempUser.Email, tempUser.Password).Scan(&tempUser).Error; err != nil {
		return nil, domain.NewInternalError("Something went wrong. Please try again later.", err)
	}
	return tempUser, nil
}

func (ur *userRepository) AddOtpData(ctx context.Context, otp *entity.Otp) (*entity.Otp, error) {

	now := time.Now()

	err := ur.db.WithContext(ctx).Raw(`
        INSERT INTO otp (email, otp, type, attempts, expires_at, created_at, deleted_at)
        VALUES (?, ?, ?, ?, ?, ?, NULL)
        ON CONFLICT (email, type) WHERE deleted_at IS NULL
        DO UPDATE SET  
            otp        = EXCLUDED.otp,
            attempts   = 0,
            expires_at = EXCLUDED.expires_at,
            updated_at = NOW(),
        RETURNING *`,
		otp.Email,
		otp.Otp,
		otp.Type,
		0,
		otp.ExpiresAt,
		now,
	).Scan(otp).Error

	if err != nil {
		ur.logger.Error("database error", "method", "add otp data", "error", err)
		return nil, domain.NewInternalError("Something went wrong. Please try again later.", err)
	}

	return otp, nil

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

func (ur *userRepository) GetLatestOtpRecord(ctx context.Context, email string, otpType entity.OtpType) (*entity.Otp, error) {
	var otpRec entity.Otp

	err := ur.db.WithContext(ctx).
		Where("email = ? AND type = ? AND deleted_at IS NULL", email, otpType).
		First(&otpRec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ur.logger.Warn("otp record not found", "email", email, "type", otpType)
			return nil, domain.NewNotFoundError("OTP record not found")
		}
		ur.logger.Error("database error", "method", "GetLatestOtpRecord", "email", email, "error", err)
		return nil, domain.NewInternalError("Something went wrong. Please try again later.", err)
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

func (ur *userRepository) CheckEmailExistInTempUser(ctx context.Context, email string) (bool, error) {

	var count int64
	if err := ur.db.WithContext(ctx).
		Where("email = ?", email).
		Model(&entity.TempUser{}).
		Count(&count).Error; err != nil {
		ur.logger.Error("database error", "error", err)
		return false, domain.NewInternalError("Something went wrong. Please try again later.", err)
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (ur *userRepository) UpdateOtpAttempts(ctx context.Context, email string, types entity.OtpType) error {

	if err := ur.db.WithContext(ctx).Where("email=? and type=?", email, types).Model(&entity.Otp{}).Update("attempts", gorm.Expr("attempts + ?", 1)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ur.logger.Warn("user email not found", "email", email, "method", "updateOtpAttempts")
			return domain.NewNotFoundError("Email not found")
		}
		ur.logger.Error("datbase error", "method", "updateOtpAttempts", "error", err)
		return domain.NewInternalError("Something went wrong. Please try again later.", err)
	}
	return nil
}

func (ur *userRepository) UpdatePassword(ctx context.Context, email string, password string) error {
	if err := ur.db.WithContext(ctx).Where("email=?", email).Model(&entity.User{}).Update("password=?", password).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ur.logger.Warn("user email not found", "email", email)
			return domain.NewNotFoundError("Email not found")
		}
		ur.logger.Error("database error", "method", "UpdatePassword", "error", err)
		return domain.NewInternalError("Something went wrong. Please try again later.", err)
	}

	return nil
}

func (ur *userRepository) DeleteOtp(ctx context.Context, id uint) error {
	if err := ur.db.WithContext(ctx).Delete(&entity.Otp{}, id).Error; err != nil {
		ur.logger.Error("database error", "error", err, "method", "Delete otp")
		return domain.NewInternalError("Something went wrong. Please try again later", err)
	}
	return nil
}

func (ur *userRepository) UpdateUserType(ctx context.Context, userId string, role string) error {

	if err := ur.db.WithContext(ctx).Where("id=?", userId).Update("user_type=?", role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ur.logger.Warn("user email not found", "email", userId)
			return domain.NewNotFoundError("user account not found")
		}
		ur.logger.Error("databser error", "error", err, "method", "UpdateOnboardingRole")
		return domain.NewInternalError("Something went wrong. plase try again later.", err)
	}
	return nil
}


