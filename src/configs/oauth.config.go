package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type oauthInstance struct {
	instance *oauth2.Config
}

func NewOAuthInstance() *oauthInstance {
	return &oauthInstance{instance: nil}
}

func (instance *oauthInstance) GetOAuthConfigInstance() (*oauth2.Config, error) {
	if instance.instance == nil {
		if err := godotenv.Load(); err != nil {
			fmt.Println("Can't load .env files! Check your .env file and try again later: " + err.Error())
			return nil, err
		} else {
			clientRedirectURL := os.Getenv("GOOGLE_CLIENT_REDIRECT_URL")
			clientId := os.Getenv("GOOGLE_CLIENT_ID")
			clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
			instance.instance = &oauth2.Config{
				ClientID:     clientId,
				ClientSecret: clientSecret,
				RedirectURL:  clientRedirectURL,
				Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
				Endpoint:     google.Endpoint,
			}
		}
	}
	return instance.instance, nil
}
