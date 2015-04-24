package main

import (
	"golang.org/x/oauth2"
	"log"
	"time"
)

const (
	exchangeLoginCaptured = "exchange login capture"
	exchangeLoginVerified = "exchange login verified"
	oauthTokenRecieved    = "oauth token recieved"
	registered            = "idle"
	syncing               = "syncing"
	syncingerror          = "syncing error"
)

type User struct {
	Username    string
	Password    string
	Token       *oauth2.Token
	Datecreated time.Time
	LastSync    time.Time
	Folderid    string
	Changekey   string
	GCalid      string
	State       string
}

type Serializer func([]User) error

var MySerializeUsers = SerializeUsers

var m map[string]User

func (u User) Save() {
	if m == nil {
		m = make(map[string]User)
	}

	log.Printf("Storing user of %s", u.Username)
	if _, ok := m[u.Username]; !ok {
		u.Datecreated = time.Now()
	}
	m[u.Username] = u
	err := MySerializeUsers(GetUsers())
	if err != nil {
		log.Println(err)
	}
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
	i := 0
	for _, value := range m {
		users[i] = value
		i++
	}
	return users
}

func DeleteUser(username string) error {
	log.Printf("Removing user of %s", username)
	delete(m, username)
	err := MySerializeUsers(GetUsers())
	if err != nil {
		log.Println(err)
	}
	return err
}
