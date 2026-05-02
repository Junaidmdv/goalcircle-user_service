package dtos

import "github.com/Junaidmdv/goalcircle-user_service/internal/domain/entity"

type RegisterRequest struct {
	FullName        string `json:"full_name" validate:"required,min=3,max=32"`
	Email           string `json:"email" validate:"required,email"`
	PhoneNum        string `json:"phone_num" validate:"required,phone"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type RegisterResponse struct {
	Email     string `json:"email"`
	PhoneNum  string `json:"phone_num"`
	OtpStatus bool   `json:"otp_status"`
}

type VerifyOtpReq struct {
	Email    string `json:"email" validate:"required,email"`
	PhoneNum string `json:"phone_num" validate:"phone"`
	Otp      string `json:"otp" validate:"required,len=6"`
}

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ResendOtpReq struct {
	Email   string         `json:"email" validate:"required,email"`
	OtpType entity.OtpType `json:"otp-type" validate:"required"`
}

type ForgotPasswordReq struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyForgotPasswordOtpReq struct {
	Email string `json:"email" validate:"required,email"`
	Otp   string `json:"otp" validate:"required,len=6"`
}

type ResetPasswordReq struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

type RenewAcccessTokenReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}  



type LogOutReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type OnboardRoleReq struct {
	UserId   string `json:"user_id" validate:"required"`
	UserRole string `json:"user_role"`
}
