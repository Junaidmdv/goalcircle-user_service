package userprofile

import (
	"bytes"
	"context"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/repository"
	uc_dtos "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
)

type userProfileUsecase struct {
	userReop  repository.UserRepository
	diskStore repository.FileStorage
	logger    logger.Logger
}

func NewUserProfileUsecase(ur repository.UserRepository, diskStore repository.FileStorage, logger logger.Logger) *userProfileUsecase {
	return &userProfileUsecase{
		userReop:  ur,
		diskStore: diskStore,
		logger:    logger,
	}
}

func (up *userProfileUsecase) GetUserProfile(ctx context.Context, input *uc_dtos.GetUserProfileReq) (*uc_dtos.GetUserProfileRes, error) {

	user, err := up.userReop.GetUserById(ctx, input.Id)
	if err != nil {
		return nil, err
	}

	return &uc_dtos.GetUserProfileRes{
		Id:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Phone:    *user.Phone,
		Street:   *user.Street,
		City:     *user.City,
		Pincode:  *user.PinCode,
		Country:  *user.Country,
		Avatar:   *user.Avatar,
	}, nil
}

func (up *userProfileUsecase) UpdateUserProfile(ctx context.Context, input *uc_dtos.UpdateUserProfileReq) (*uc_dtos.UpdateUserProfileRes, error) {

	updates := map[string]interface{}{}

	if input.Phone != "" {
		updates["phone"] = input.Phone
	}
	if input.Street != "" {
		updates["street"] = input.Street
	}
	if input.City != "" {
		updates["city"] = input.City
	}
	if input.State != "" {
		updates["state"] = input.State
	}
	if input.Pincode != "" {
		updates["zip"] = input.Pincode
	}
	if input.Country != "" {
		updates["country"] = input.Country
	}

	if len(updates) == 0 {
		return nil, domain.NewBadRequestError("no fields to update")
	}
	user, err := up.userReop.UpdateUser(ctx, input.Id, updates)
	if err != nil {
		return nil, err
	}
	return &uc_dtos.UpdateUserProfileRes{
		Id:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Phone:    *user.Phone,
		Street:   *user.Street,
		City:     *user.City,
		State:    *user.State,
		Pincode:  *user.PinCode,
		Country:  *user.Country,
	}, nil

}

func (up *userProfileUsecase) UpdateUserProfileImage(ctx context.Context, input *uc_dtos.UpdateUserProfileImgReq) (*uc_dtos.UpdateUserProfileImgRes, error) {

	url, err := up.diskStore.UploadFile(ctx, repository.Logo, bytes.NewReader(input.Data), &repository.FileMetadata{
		Filename:    input.FileName,
		ContentType: input.MimeType,
		Size:        int64(len(input.Data)),
	})
	if err != nil {
		up.logger.Error("upload file failure", "method", "UpdateUserProfileImage", "error", err)
		return nil, domain.NewInternalError("Something went wrong.Please try again later", err)
	}

	if err := up.userReop.UpdateUserProfileImage(ctx, input.Id, url); err != nil {
		return nil, err
	}

	return &uc_dtos.UpdateUserProfileImgRes{
		Url: url,
	}, nil
}

func (up *userProfileUsecase) ChangePassword(ctx context.Context, input *uc_dtos.ChangePasswordReq) (*uc_dtos.ChangePasswordRes, error) {  
	    
	

	return nil, nil
}
