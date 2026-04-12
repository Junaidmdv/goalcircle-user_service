package dtos

type RegisterRequest struct {
	FullName        string `validate:"required,min=3,max=32"`
	Email           string `validate:"required,email"` 
	PhoneNum        string `validate:"required,phone"`
	Password        string `validate:"required"`
	ConfirmPassword string `validate:"required"`
}

type RegisterResponse struct {
	UserId    string `json:"user_id"`
	Email     string `json:"email"`
	OtpStatus bool   `json:"otp_status"`
}
