package oauth

import (
	"time"

	"github.com/Junaidmdv/goalcircle-user_service/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOauth struct {
	Config  *oauth2.Config
	TimeOut time.Duration
}

func NewGoogleOauth(cnfg *config.GoogleAuthConfig) *GoogleOauth {
	googleConfig := &oauth2.Config{
		ClientID:     cnfg.ClientId,
		ClientSecret: cnfg.ClientSecret,
		RedirectURL:  cnfg.RedirectUrl,
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}

	return &GoogleOauth{
		Config:  googleConfig,
		TimeOut: cnfg.TimeOut,
	}
}
