package entity 


type OtpType string 


const (
	 Register OtpType="register" 
	 ResetPassword OtpType="reset_pasword" 
	 ResendOtp OtpType="resended_otp" 
	 ForgotPassword OtpType="forgot password"  
	 OtpMaxAttempts=5

)   