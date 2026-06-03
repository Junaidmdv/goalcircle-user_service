package usermanagement

import (
	"context"

	uc_dtos "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"
)

type UserManagementUsecase interface {
	BlockUser(context.Context, *uc_dtos.BlockUserReq) (*uc_dtos.BlockUserRes, error)
	UnBlockUser(context.Context, *uc_dtos.UnblockUserReq) (*uc_dtos.UnblockUserRes, error) 
	GetUsers(context.Context)([]*uc_dtos.GetUserRes,error)
}


