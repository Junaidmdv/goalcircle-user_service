package userprofile

import (
	"context"

	uc_dtos "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"
)

type UserProfileUsecase interface {
	GetUserProfile(context.Context, *uc_dtos.GetUserProfileReq) (*uc_dtos.GetUserProfileRes, error)
	UpdateUserProfile(context.Context, *uc_dtos.UpdateUserProfileReq) (*uc_dtos.UpdateUserProfileRes, error)
	UpdateUserProfileImage(context.Context, *uc_dtos.UpdateUserProfileImgReq) (*uc_dtos.UpdateUserProfileImgRes, error)
	ChangePassword(context.Context, *uc_dtos.ChangePasswordReq) (*uc_dtos.ChangePasswordRes, error)
}
