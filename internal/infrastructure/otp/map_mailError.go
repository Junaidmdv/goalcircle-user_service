package otp

import (
	"strings"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain"
)

func (e *EmailService) MapMailError(err error) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	switch {
	// Network / dial errors
	case strings.Contains(errMsg, "no such host"),
		strings.Contains(errMsg, "connection refused"):
		return domain.NewInternalError("email server unreachable", err)

	case strings.Contains(errMsg, "x509"),
		strings.Contains(errMsg, "tls"):
		return domain.NewInternalError("email server TLS error", err)

	// Auth errors
	case strings.Contains(errMsg, "535"),
		strings.Contains(errMsg, "534"),
		strings.Contains(errMsg, "530"),
		strings.Contains(errMsg, "454"):
		return domain.NewUnAuthenticatedError("email service authentication failed")

	// Invalid recipient — user input problem
	case strings.Contains(errMsg, "550"),
		strings.Contains(errMsg, "551"),
		strings.Contains(errMsg, "553"),
		strings.Contains(errMsg, "501"):
		return domain.NewValidationError("recipient email address is invalid or does not exist")

	// Mailbox full — not your fault, not user's fault
	case strings.Contains(errMsg, "552"):
		return domain.NewInternalError("recipient mailbox is full", err)

	// Rate limit
	case strings.Contains(errMsg, "421"):
		return domain.NewDeadlineExceedeError("email service rate limit reached")

	default:
		return domain.NewInternalError("failed to send otp email", err)
	}
}
