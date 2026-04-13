package twilio

import (
	"fmt"
	"time"

	"github.com/junaidmdv/goalcircle/user_service/internal/config"
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

type smsOtpService struct {
	client *twilio.RestClient
	params *SMSParams
}

func NewSMSOtpService(config *config.TwilioConfig) OtpService {
	cl := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.AccountSID,
		Password: config.AuthToken,
	})

	p := &SMSParams{
		FromNum:    config.FromNum,
		ExpireTime: time.Duration(config.OtpExpiryTime),
	}

	return &smsOtpService{
		client: cl,
		params: p,
	}
}

type OtpRes struct {
	Otp       string
	ExpiresAt time.Time
}

func (s *smsOtpService) SendOtp(num string) (*OtpRes, error) {
	s.params.TONum = num
	if err := s.params.GenerateOtp(6); err != nil {
		return nil, err
	}
	body := fmt.Sprintf("Your OTP is: %s. Valid for 5 minutes. Do not share it.", s.params.Otp)

	params := &openapi.CreateMessageParams{}
	params.SetFrom(s.params.FromNum)
	params.SetTo(s.params.TONum)
	params.SetBody(body)
	_, err := s.client.Api.CreateMessage(params)
	if err != nil {
		return nil, err
	}

	return &OtpRes{
		Otp:       s.params.Otp,
		ExpiresAt: time.Now().Add(s.params.ExpireTime),
	}, nil

}
