package twilio

import (
	"crypto/rand"
	"math"
	"math/big"
	"time"
)

type SMSParams struct {
	FromNum   string
	TONum     string
	Otp       string
	ExpiresAt time.Duration
}

func (s *SMSParams) GenerateOtp(maxdigit uint32) error {
	bi, err := rand.Int(
		rand.Reader,
		big.NewInt(int64(math.Pow(10, float64(maxdigit)))),
	)
	if err != nil {
		return err
	}
	s.Otp = bi.String()
	return nil
}
