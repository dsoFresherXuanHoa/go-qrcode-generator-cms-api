package configs

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	ErrLoadOAuthEnvFile = errors.New("load .env file failure")
)

type oauthClient struct {
	instance *oauth2.Config
}

func NewOAuthClient() *oauthClient {
	return &oauthClient{instance: nil}
}

func (instance *oauthClient) Instance() (*oauth2.Config, error) {
	if instance.instance == nil {
		if err := godotenv.Load(); err != nil {
			fmt.Println("Error while load .env file: " + err.Error())
			return nil, ErrLoadOAuthEnvFile
		} else {
			clientRedirectURL := os.Getenv("GOOGLE_CLIENT_REDIRECT_URL")
			clientId := os.Getenv("GOOGLE_CLIENT_ID")
			clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
			clientScopes := []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"}
			instance.instance = &oauth2.Config{
				ClientID:     clientId,
				ClientSecret: clientSecret,
				RedirectURL:  clientRedirectURL,
				Scopes:       clientScopes,
				Endpoint:     google.Endpoint,
			}
		}
	}
	return instance.instance, nil
}
