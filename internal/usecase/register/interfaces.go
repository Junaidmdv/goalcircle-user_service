package register

import (
	"context"
	uc_dtos "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"
)

type RegistrationUsecase interface {
	InitiateUserRegistration(context.Context, *uc_dtos.RegisterRequest) (*uc_dtos.RegisterResponse, error)
}
