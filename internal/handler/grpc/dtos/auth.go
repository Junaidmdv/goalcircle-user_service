package dtos

type RegisterRequest struct {
	FullName        string `json:"full_name" validate:"required,min=3,max=32"`
	Email           string `json:"email" validate:"required,email"` 
	PhoneNum        string `json:"phone_num" validate:"required,phone"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

type RegisterResponse struct {
	UserId    string `json:"user_id"`
	Email     string `json:"email"`
	OtpStatus bool   `json:"otp_status"`
}
