package main

import (
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type OAuthConfig interface {
	ClientId() string
	ClientSecret() string
	RedirectURL() string
}

var conf = &oauth2.Config{
	Scopes:   []string{calendar.CalendarScope},
	Endpoint: google.Endpoint,
}

func SetConfig(config OAuthConfig) {
	conf.ClientID = config.ClientId()
	conf.ClientSecret = config.ClientSecret()
	conf.RedirectURL = config.RedirectURL()
}

func getClient(user User) *http.Client {
	return conf.Client(oauth2.NoContext, user.Token)
}

func tryOAuth2(w http.ResponseWriter, r *http.Request, user User) {
	log.Println("Starting OAuth2")

	url := conf.AuthCodeURL(user.Username)
	log.Println("Going to %v", url)
	http.Redirect(w, r, string(url), 302)
}

func handleOAuthResponse(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("state")
	log.Printf("Username is: %s", username)

	code := r.URL.Query().Get("code")
	log.Printf("Code is: %s", code)
	// Exchanging the code for a token
	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Printf("Token is %s", token)

	user := GetUser(username)
	user.Token = token
	user.State = registered
	user.Save()
}
