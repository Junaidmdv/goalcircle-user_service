package dtos

import "time"

type RegisterRequest struct {
	FullName        string
	Email           string
	PhoneNum        string
	Password        string
	ConfirmPassword string
}

type RegisterResponse struct {
	Email     string
	PhoneNum  string
	OtpStatus bool
	OtpExpiry time.Time
}

type VerifyOtpRequest struct {
	Email    string
	PhoneNum string
	Otp      string
}

type VerifyOtpResponse struct {
	UserId            string
	AccessToken       string
	AceessTokenExpiry time.Time
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	UserId            string
	AccessToken       string
	AccessTokenExpiry time.Time
}

type ResendOtpReq struct {
	Email    string
	PhoneNum string
}

type ResendOtpResponse struct {
	Success bool 
	OtpExpiry time.Time 
}


