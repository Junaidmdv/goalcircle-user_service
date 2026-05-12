package oauth

import (
	"github.com/Junaidmdv/goalcircle-user_service/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOuth struct {
	config *oauth2.Config
}

func NewGoogleOauth(cnfg *config.GoogleAuthConfig) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cnfg.ClientId,
		ClientSecret: cnfg.ClientSecret,
		RedirectURL:  cnfg.RedirectUrl,
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}
}
