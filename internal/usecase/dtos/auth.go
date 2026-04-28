package dtos

import (
	"time"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/entity"
)

type RegisterRequest struct {
	FullName        string
	Email           string
	Password        string
	ConfirmPassword string
}

type RegisterResponse struct {
	Email     string
	OtpStatus bool
	OtpExpiry time.Time
}

type VerifyOtpRequest struct {
	Email string
	Otp   string
}

type VerifyOtpResponse struct {
	Success            bool
	SessionId          string
	UserId             string
	FullName           string
	Email              string
	AccessToken        string
	AceessTokenExpiry  time.Time
	RefreshToken       string
	RefreshTokenExpiry time.Time
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	SessionId          string
	UserId             string
	Email              string
	FullName           string
	AccessToken        string
	AccessTokenExpiry  time.Time
	RefreshToken       string
	RefreshTokenExpiry time.Time
}

type ResendOtpReq struct {
	Email   string
	OtpType entity.OtpType
}

type ResendOtpResponse struct {
	Success   bool
	OtpExpiry time.Time
}

type ForgotPasswordReq struct {
	Email string
}

type ForgotPasswordRes struct {
	Success   bool
	ExpiresAt time.Time
}

type VerifyForgotPasswordOtpReq struct {
	Email string
	Otp   string
}

type VerifyForgotPasswordOtpRes struct {
	Success    bool
	ResetToken string
	ExpiresAt  time.Time
}

type ResetPasswordRes struct {
	Success bool
}

type ResetPasswordReq struct {
	Email      string
	Password   string 
	ResetToken string
}

type RenewAcccessTokenReq struct {
	RefreshToken string
}

type RenewAccessTokenRes struct {
	AccessToken       string
	AccessTokenExpiry time.Time
}

type LogOutReq struct {
	RefreshToken string
}

type LogOutRes struct {
	Success bool
}

type OnboardingRoleReq struct {
	UserId string
	Role   string
}

type OnboardingRoleRes struct {
	Success bool
}

type OnboardingTeamDtlsReq struct {
}

type OnboardingTeamDtlsRes struct {
}

type OnboardingOrganiserDtlsReq struct {
}

type OnboardingAddOrganiserDtlsRes struct {
}
