package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
)

var GoogleOAuthConfig *oauth2.Config

func InitGoogleAuth() {

	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("Client_ID"),
		ClientSecret: os.Getenv("Client_Secret"),
		RedirectURL:  "",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
