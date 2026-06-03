package dtos 

type AdminRegisterRequest struct {
	FullName        string `json:"full_name" validate:"required,min=3,max=32"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,password"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

// type AdminRegisterResponse struct {
// 	Email     string `json:"email"`
// 	OtpStatus bool   `json:"otp_status"`
// }

type AdminLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}
