package dtos

type RegisterRequest struct {
	FullName        string
	Email           string
	Password        string
	ConfirmPassword string
}

type RegisterResponse struct {
	SessionId    string
	Email     string
	OtpStatus bool
}


