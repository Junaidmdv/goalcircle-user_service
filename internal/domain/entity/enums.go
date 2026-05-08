package entity

type OtpType string

const (
	Register       OtpType = "register"
	ResetPassword  OtpType = "reset_pasword"
	ForgotPassword OtpType = "forgot_password" 
	OtpMaxAttempts         = 5
)


const (
	UNSPECIFIED = "unspecified"
	ORGAINISER  = "organiser"
	MANAGER     = "manager"
)
