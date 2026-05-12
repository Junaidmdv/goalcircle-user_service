package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Junaidmdv/goalcircle-user_service/internal/config"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain"
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

func (gl *GoogleOauth) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	goauth := gl.Config

	token, err := goauth.Exchange(ctx, code)
	if err != nil {
		var retrieveErr *oauth2.RetrieveError
		if errors.As(err, &retrieveErr) {
			switch retrieveErr.Response.StatusCode {
			case 400:
				return nil, domain.NewBadRequestError("Authorization code is invalid or expired")
			case 401:
				return nil, domain.NewUnAuthenticatedError("OAuth authentication failed")
			case 500, 503:
				return nil, domain.NewInternalError("Google service unavailable", err)
			default:
				return nil, domain.NewInternalError("Something went wrong", err)
			}
		}
		return nil, domain.NewInternalError("Something went wrong", err)
	}
	return token, nil
}

type GoogleUserInfo struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	Verified bool   `json:"verified_email"`
}

func (gl *GoogleOauth) GetUserData(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {

	contxt, cancle := context.WithTimeout(ctx, time.Second*5)
	defer cancle()

	client := gl.Config.Client(contxt, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, domain.NewInternalError("Something went wrong. Please try again later.", err)
	}
	defer resp.Body.Close()

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, domain.NewInternalError("Something went wrong", fmt.Errorf("failed to unmarshall user data"))
	}
	return &userInfo, nil
}
