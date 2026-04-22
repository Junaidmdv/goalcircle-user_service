package entity 


type OtpType string 


const (
	 Register OtpType="register" 
	 ResetPassword OtpType="reset_pasword" 
	 ResendOtp OtpType="resended_otp"
)