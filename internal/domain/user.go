package domain


type TempUser struct {
    Name      string `json:"name"`
    Email     string `json:"email"`
    Phone     string `json:"phone"`
    Password  string `json:"password_hash"` 
    OtpHash   string `json:"otp_hash"`       // SHA-256 of plaintext OTP
    Attempts  int    `json:"attempts"`
    CreatedAt int64  `json:"created_at"`
}