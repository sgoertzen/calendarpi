package main

import (
	"golang.org/x/oauth2"
	"log"
	"time"
)

type User struct {
	Username    string
	Password    string
	Token       *oauth2.Token
	Datecreated time.Time
	Folderid    string
	Changekey   string
	State       int
}

const (
	exchangeLoginCaptured = 1
	exchangeLoginVerified = 20
	oauthTokenRecieved    = 50
	registered            = 100
)

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
	err := SerializeUsers(GetUsers())
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

func DeleteUser(username string) {
	log.Printf("Removing user of %s", username)
	delete(m, username)
	err := SerializeUsers(GetUsers())
	if err != nil {
		log.Println(err)
	}
}
