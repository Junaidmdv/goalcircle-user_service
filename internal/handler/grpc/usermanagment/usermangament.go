package usermanagement

import (
	"context"
	"time"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain"
	"github.com/Junaidmdv/goalcircle-user_service/internal/handler/grpc/auth"
	dt "github.com/Junaidmdv/goalcircle-user_service/internal/handler/grpc/dtos"
	"github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"
	usermanagementUc "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/usermanagement"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/validater"
	"github.com/Junaidmdv/goalcircle-user_service/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type adminUserManagement struct {
	userManagementUc usermanagementUc.UserManagementUsecase
	timeout          *time.Duration
	logger           logger.Logger
	validater        *validater.Validater
	pb.UnimplementedAdminUserManagementServiceServer
}

func NewAdminUserManagementHandler(mu usermanagementUc.UserManagementUsecase, timeout *time.Duration, logger logger.Logger, validater validater.Validater) *adminUserManagement {
	return &adminUserManagement{
		userManagementUc: mu,
		timeout:          timeout,
		logger:           logger,
		validater:        &validater,
	}
}

func (um *adminUserManagement) BlockUser(ctx context.Context, input *pb.BlockUserReq) (*pb.BlockUserRes, error) {
	context, cancel := context.WithTimeout(ctx, *um.timeout)
	defer cancel()

	if validationErrs := um.validater.Validation(&dt.BlockUserReq{
		UserId: input.UserId,
	}); validationErrs != nil {
		stWithDetails, err := auth.ValidationError(validationErrs)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	res, err := um.userManagementUc.BlockUser(context, &dtos.BlockUserReq{
		UserId: input.UserId,
	})
	if err != nil {
		return nil, domain.GRPCStatus(err)
	}

	return &pb.BlockUserRes{
		Success: res.Success,
	}, nil
}

func (um *adminUserManagement) UnBlockUser(ctx context.Context, input *pb.UnblockUserReq) (*pb.UnblockUserRes, error) {
	context, cancel := context.WithTimeout(ctx, *um.timeout)
	defer cancel()

	if validationErrs := um.validater.Validation(&dt.UnBlockUserReq{
		UserId: input.UserId,
	}); validationErrs != nil {
		stWithDetails, err := auth.ValidationError(validationErrs)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	res, err := um.userManagementUc.UnBlockUser(context, &dtos.UnblockUserReq{
		UserId: input.UserId,
	})
	if err != nil {
		return nil, domain.GRPCStatus(err)
	}

	return &pb.UnblockUserRes{
		Success: res.Success,
	}, nil
}

func (um *adminUserManagement) GetUsers(ctx context.Context, res *emptypb.Empty) (*pb.GetUserResponse, error) {

	users, err := um.userManagementUc.GetUsers(ctx)

	if err != nil {
		return nil, domain.GRPCStatus(err)
	}
	var userList []*pb.User

	for _, user := range users {
		userList = append(userList, &pb.User{
			UserId:    user.UserId,
			Email:     user.Email,
			IsBlocked: user.IsBlocked,
			CreatedAt: timestamppb.New(user.CreatedAt),
		})
	}

	return &pb.GetUserResponse{
		Users: userList,
	}, nil

}
