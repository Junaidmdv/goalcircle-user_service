package adminauth

import (
	"context"

	uc_dtos "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"
)

type AdminAuthUsecase interface {
	Register(context.Context, *uc_dtos.AdminAuthRegisterReq) (*uc_dtos.AdminAuthRegisterRes, error)
	Login(context.Context, *uc_dtos.AdminLoginReq) (*uc_dtos.AdminLoginRes, error)
}
