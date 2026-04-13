package usecase

import (
	"context"

	"github.com/junaidmdv/goalcircle/user_service/internal/usecase/dtos"
)

type AuthUsecase interface {
	InitiateUserRegistration(context.Context, *dtos.RegisterRequest) (*dtos.RegisterResponse, error)
}
