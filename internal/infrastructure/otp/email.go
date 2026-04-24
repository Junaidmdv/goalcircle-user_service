package otp

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/Junaidmdv/goalcircle-user_service/internal/config"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	server gomail.Dialer
	config *config.SMTPConfig
}

type OtpResponse struct {
	Otp    string
	Expiry time.Duration
}

// Constructor
func NewEmailService(cfg *config.SMTPConfig) (*EmailService, error) {

	dialer := gomail.NewDialer(cfg.ServerURL, cfg.Port, cfg.FromEmail, cfg.Password)

	return &EmailService{
		server: *dialer,
		config: cfg,
	}, nil
}

// Generate secure 6-digit OTP
func (e *EmailService) GenerateOTP(length int) (string, error) {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(length)), nil)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%0*d", length, n), nil

}

// Send OTP email
func (e *EmailService) SendOTP(toEmail string) (*OtpResponse, error) {
	otp, err := e.GenerateOTP(6)
	if err != nil {
		return nil, err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", e.server.Username)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Your OTP Verification Code")
	m.SetBody("text/html", fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="margin:0;padding:0;background-color:#f4f6f9;font-family:Arial,sans-serif;">

  <table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f4f6f9;padding:40px 0;">
    <tr>
      <td align="center">
        <table width="480" cellpadding="0" cellspacing="0" style="background-color:#ffffff;border-radius:12px;overflow:hidden;border:1px solid #e8eaed;">

          <!-- Header -->
          <tr>
            <td style="background-color:#185FA5;padding:32px;text-align:center;">
              <h1 style="margin:0;color:#ffffff;font-size:22px;font-weight:600;letter-spacing:0.5px;">
                Password Reset
              </h1>
              <p style="margin:8px 0 0;color:#b3d4f0;font-size:13px;">
                One-Time Password Verification
              </p>
            </td>
          </tr>

          <!-- Body -->
          <tr>
            <td style="padding:40px 32px 24px;">
              <p style="margin:0 0 8px;color:#1a1a2e;font-size:16px;font-weight:600;">
                Hello,
              </p>
              <p style="margin:0 0 28px;color:#555770;font-size:14px;line-height:1.6;">
                We received a request to reset your password. Use the OTP below to proceed.
                This code is valid for <strong style="color:#185FA5;">10 minutes</strong> and can only be used once.
              </p>

              <!-- OTP Box -->
              <table width="100%%" cellpadding="0" cellspacing="0">
                <tr>
                  <td align="center">
                    <div style="background-color:#f0f6ff;border:1.5px dashed #185FA5;border-radius:12px;padding:28px 0;margin-bottom:28px;">
                      <p style="margin:0 0 8px;color:#888;font-size:12px;letter-spacing:1px;text-transform:uppercase;">
                        Your OTP Code 
                      </p>
                      <p style="margin:0;font-size:42px;font-weight:700;letter-spacing:16px;color:#185FA5;padding-left:16px;">
                        %s
                      </p>
                    </div>
                  </td>
                </tr>
              </table>

              <!-- Warning -->
              <table width="100%%" cellpadding="0" cellspacing="0">
                <tr>
                  <td style="background-color:#fff8f0;border-left:4px solid #EF9F27;border-radius:0 8px 8px 0;padding:12px 16px;margin-bottom:24px;">
                    <p style="margin:0;color:#854F0B;font-size:13px;line-height:1.5;">
                      <strong>Security Notice:</strong> Never share this OTP with anyone.
                      Our team will never ask for your OTP.
                    </p>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- Divider -->
          <tr>
            <td style="padding:0 32px;">
              <hr style="border:none;border-top:1px solid #f0f0f0;margin:0;">
            </td>
          </tr>

          <!-- Footer -->
          <tr>
            <td style="padding:24px 32px;text-align:center;">
              <p style="margin:0 0 4px;color:#aaa;font-size:12px;">
                If you didn't request a password reset, please ignore this email.
              </p>
              <p style="margin:0;color:#aaa;font-size:12px;">
                This OTP will expire automatically after 10 minutes.
              </p>
            </td>
          </tr>

        </table>
      </td>
    </tr>
  </table>

</body>
</html>
	`, otp))

	if err := e.server.DialAndSend(m); err != nil {
		return nil, err
	}

	return &OtpResponse{
     Otp: otp, 
     Expiry: e.config.OtpExpiry,
  },nil
}

// Verify OTP
func (e *EmailService) VerifyOTP(storedOTP, storedExpiry, inputOTP string) error {
	// Check expiry
	expiry, err := time.Parse(time.RFC3339, storedExpiry)
	if err != nil {
		return fmt.Errorf("invalid expiry format: %w", err)
	}

	if time.Now().After(expiry) {
		return fmt.Errorf("otp has expired")
	}

	// Check match
	if storedOTP != inputOTP {
		return fmt.Errorf("invalid otp")
	}

	return nil
}
