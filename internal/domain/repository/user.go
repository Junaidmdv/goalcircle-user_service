package repository

import (
	"context"

	"github.com/junaidmdv/goalcirlcle/user_service/internal/domain/entity"
)

type UserRepository interface {
	ExistByEmail(context.Context, string) (bool, error)
	CreateTempUser(context.Context, *entity.TempUser) (*entity.TempUser, error)
}
