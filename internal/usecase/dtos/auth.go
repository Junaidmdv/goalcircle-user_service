package dtos

type RegisterRequest struct {
	FullName        string
	Email           string
	PhoneNum        string
	Password        string
	ConfirmPassword string
}

type RegisterResponse struct {
	UserId    string 
	Email     string
	OtpStatus bool
}
