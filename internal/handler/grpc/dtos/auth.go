package dtos

import "time"

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

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	UserId            string
	AccessToken       string
	AccessTokenExpiry time.Time
}
