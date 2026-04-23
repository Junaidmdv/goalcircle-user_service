package otp

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/Junaidmdv/goalcircle-user_service/internal/config"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type OtpService interface {
	SendOtp(string) (*OtpRes, error)
	ParseTwilioError(error) error
}

type ChannelType int

const (
	Email ChannelType = iota
	SMS
	WhatsApp
)

type SmsOtpService struct {
	client     *twilio.RestClient
	FromNum    string
	ExpireTime time.Duration
}

func NewSMSOtpService(config *config.TwilioConfig) OtpService {
	cl := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.AccountSID,
		Password: config.AuthToken,
	})

	return &SmsOtpService{
		client:     cl,
		FromNum:    config.FromNum,
		ExpireTime: time.Duration(config.OtpExpiryTime),
	}
}

type OtpRes struct {
	Otp       string
	ExpiresAt time.Time
}

func (s *SmsOtpService) SendOtp(num string) (*OtpRes, error) { 
	fmt.Println(num)

	otp, err := generateRandomOtp(6)
	if err != nil {
		return nil, err
	}
	body := fmt.Sprintf("Your OTP is: %s. Valid for 5 minutes. Do not share it.", otp)

	params := &openapi.CreateMessageParams{}
	params.SetFrom(s.FromNum)
	params.SetTo(num)
	params.SetBody(body)
	res, err := s.client.Api.CreateMessage(params)
	if err != nil {
		return nil, err
	}
	log.Println(res)

	return &OtpRes{
		Otp:       otp,
		ExpiresAt: time.Now().Add(s.ExpireTime),
	}, nil

}

func generateRandomOtp(length int) (string, error) {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(length)), nil)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%0*d", length, n), nil
}
