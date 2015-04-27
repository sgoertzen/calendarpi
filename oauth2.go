package main

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type OAuthConfig interface {
	ClientId() string
	ClientSecret() string
	RedirectURL() string
	MaxFetchSize() int
}

var maxFetchSize int

var conf = &oauth2.Config{
	Scopes:   []string{calendar.CalendarScope},
	Endpoint: google.Endpoint,
}

func SetOauthConfig(config OAuthConfig) {
	conf.ClientID = config.ClientId()
	conf.ClientSecret = config.ClientSecret()
	conf.RedirectURL = config.RedirectURL()
	maxFetchSize = config.MaxFetchSize()
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

func handleOAuthResponse(w http.ResponseWriter, r *http.Request) (User, error) {
	username := r.URL.Query().Get("state")
	log.Printf("Username is: %s", username)

	code := r.URL.Query().Get("code")
	log.Printf("Code is: %s", code)
	// Exchanging the code for a token
	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return User{}, err
	}
	log.Printf("Token is %s", token)

	user := GetUser(username)
	user.Token = token
	user.State = registered
	user.Save()
	return user, nil
}

func GetGCalendarList(user User) *calendar.CalendarList {
	client := getClient(user)
	srv, err := calendar.New(client)
	if err != nil {
		log.Println("Unable to get calendar service", err)
	}
	calendars, err2 := srv.CalendarList.List().MinAccessRole("writer").Do()
	if err2 != nil {
		log.Println("Unable to get calendar list", err)
	}
	return calendars
}

func getGCalAppointments(user User) (*calendar.Events, error) {
	client := getClient(user)
	srv, err := calendar.New(client)
	if err != nil {
		log.Println("Unable to retrieve calendar Client %v", err)
	}
	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List(user.GCalid).ShowDeleted(false).
		SingleEvents(true).
		TimeMin(t).MaxResults(int64(maxFetchSize)).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve the user's events. %v", err)
	}
	return events, nil
}
