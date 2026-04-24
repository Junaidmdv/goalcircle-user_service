package otp

import "gopkg.in/gomail.v2"

type EmailOtpService struct {
	server gomail.Dialer
} 


package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
	"your_project/config"
)

type EmailOtpService struct {
	server gomail.Dialer
}

// Constructor
func NewEmailOtpService(cfg *config.SMTP) (*EmailOtpService, error) {
	port, err := strconv.Atoi(cfg.PORT)
	if err != nil {
		return nil, fmt.Errorf("invalid smtp port: %w", err)
	}

	dialer := gomail.NewDialer(cfg.ServerURL, port, cfg.FromEmail, cfg.Password)

	return &EmailOtpService{
		server: *dialer,
	}, nil
}

// Generate secure 6-digit OTP
func (e *EmailOtpService) GenerateOTP() (string, error) {
	max := big.NewInt(999999)
	min := big.NewInt(100000)

	n, err := rand.Int(rand.Reader, max.Sub(max, min))
	if err != nil {
		return "", fmt.Errorf("failed to generate otp: %w", err)
	}

	otp := n.Add(n, min).String()
	return otp, nil
}

// Send OTP email
func (e *EmailOtpService) SendOTP(toEmail string) (string, error) {
	otp, err := e.GenerateOTP()
	if err != nil {
		return "", err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", e.server.Username)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Your OTP for Password Reset")
	m.SetBody("text/html", fmt.Sprintf(`
		<div style="font-family:Arial,sans-serif;max-width:480px;margin:auto;padding:32px;border:1px solid #eee;border-radius:8px;">
			<h2 style="color:#185FA5;margin-bottom:8px;">Password Reset OTP</h2>
			<p style="color:#555;">Use the OTP below to reset your password. It expires in <strong>10 minutes</strong>.</p>
			<div style="text-align:center;margin:32px 0;">
				<span style="font-size:36px;font-weight:bold;letter-spacing:12px;color:#185FA5;">%s</span>
			</div>
			<p style="color:#999;font-size:12px;">If you didn't request this, please ignore this email.</p>
		</div>
	`, otp))

	if err := e.server.DialAndSend(m); err != nil {
		return "", fmt.Errorf("failed to send otp email: %w", err)
	}

	return otp, nil
}

// Verify OTP
func (e *EmailOtpService) VerifyOTP(storedOTP, storedExpiry, inputOTP string) error {
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