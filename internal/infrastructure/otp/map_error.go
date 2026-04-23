package otp

import (
	"errors"
	"fmt"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain"
	"github.com/twilio/twilio-go/client"
)

func (s *SmsOtpService) ParseTwilioError(err error) error {
	// cast to Twilio's error type
	var twilioErr *client.TwilioRestError
	if errors.As(err, &twilioErr) {
		switch twilioErr.Code {
		// user errors — invalid number, unverified etc
		case 21211:
			return domain.NewValidationError("invalid phone number format")
		case 21608:
			return domain.NewValidationError("phone number is not verified on free tier")
		case 21610:
			return domain.NewValidationError("phone number is blacklisted")
		case 21614:
			return domain.NewValidationError("phone number is not SMS capable")

		// server/config errors
		case 20003:
			return domain.NewInternalError("invalid Twilio credentials", errors.New(twilioErr.Message))
		case 21212:
			return domain.NewInternalError("invalid From number", errors.New(twilioErr.Message))
		case 30008:
			return domain.NewInternalError("message delivery failed", errors.New(twilioErr.Message))

		default:
			return domain.NewInternalError(fmt.Sprintf("twilio error %d: %s", twilioErr.Code, twilioErr.Message), errors.New(twilioErr.MoreInfo))
		}
	}

	return domain.NewInternalError("failed to send OTP: %w", err)
}
