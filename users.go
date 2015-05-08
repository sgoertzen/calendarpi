package main

import (
	"github.com/sgoertzen/xchango"
	"golang.org/x/oauth2"
	"log"
	"sort"
	"time"
)

const (
	exchangeLoginCaptured = "exchange login capture"
	exchangeLoginVerified = "exchange login verified"
	oauthTokenRecieved    = "oauth token recieved"
	registered            = "queued"
	syncing               = "syncing"
	syncingerror          = "sync error"
	successfulsync        = "success"
	registererror         = "registration error"
)

type User struct {
	// Todo: Store username in both places for now.
	// TODO: Delete the username, password, Folderid and changekey once they are all moved over
	Username    string
	Password    string
	Token       *oauth2.Token
	Datecreated time.Time
	LastSync    time.Time
	Folderid    string
	Changekey   string
	GCalid      string
	State       string
	ExUser      *xchango.ExchangeUser
	ExCal       *xchango.ExchangeCalendar
}

type Serializer func([]User) error

var serializeUsers = SerializeUsers
var m map[string]User

func (u User) Save() error {
	if m == nil {
		m = make(map[string]User)
	}
	// These are here temporarily to populate our new objects.  We can remove this once all users are migrated
	if u.ExUser == nil {
		u.ExUser = &xchango.ExchangeUser{Username: u.Username, Password: u.Password}
	}
	if u.ExCal == nil {
		u.ExCal = &xchango.ExchangeCalendar{Folderid: u.Folderid, Changekey: u.Changekey}
	}

	log.Printf("Storing user: '%s'", u.ExUser.Username)
	t := time.Time{}
	if u.Datecreated == t {
		u.Datecreated = time.Now()
	}
	m[u.Username] = u
	err := serializeUsers(GetUsers())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetUser(username string) User {
	if m == nil {
		return User{}
	}
	return m[username]
}

func GetUsers() []User {
	if m == nil {
		return make([]User, 0)
	}
	users := make([]User, len(m))

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, username := range keys {
		users[i] = m[username]
	}
	return users
}

func DeleteUser(username string) error {
	log.Printf("Removing user: '%s'", username)
	delete(m, username)
	err := serializeUsers(GetUsers())
	if err != nil {
		log.Println(err)
	}
	return err
}
