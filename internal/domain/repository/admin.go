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

type AdminRepository interface {
	Create(context.Context, *entity.Admin) (*entity.Admin, error)
	GetAdminByEmail(context.Context, string) (*entity.Admin, error)
}

type adminRepository struct {
	db      *gorm.DB
	logger  logger.Logger
	timeout time.Duration
}

func NewAdminRepository(db *gorm.DB, logger logger.Logger, timeout time.Duration) AdminRepository {
	return &adminRepository{
		db:      db,
		logger:  logger,
		timeout: timeout,
	}
}

func (ar *adminRepository) Create(ctx context.Context, admin *entity.Admin) (*entity.Admin, error) {
	if err := ar.db.WithContext(ctx).Create(admin).Error; err != nil {
		ar.logger.Error("database error", "method", "AdminRepository.Create")
		return nil, domain.NewInternalError("Somthing went wrong.Please try again later.", err)
	}
	return admin, nil
}

func (ar *adminRepository) GetAdminByEmail(ctx context.Context, email string) (*entity.Admin, error) {
	var user entity.Admin
	if err := ar.db.WithContext(ctx).
		Where("email=?", email).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ar.logger.Warn("user email not found", "email", email)
			return nil, domain.NewNotFoundError("No account found with this email. Please register first.")
		}
		ar.logger.Error("database error", "error", err)
		return nil, domain.NewInternalError("Something went wrong. Please try again later.", err)

	}
	return &user, nil
}


