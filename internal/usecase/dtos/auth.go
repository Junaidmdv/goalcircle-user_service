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

type OtpRequest struct {
	Email    string
	PhoneNum string
	Otp      string
}



type OtpResponse struct{
  
}