package main

import (
	"github.com/stretchr/testify/assert"
	"errors"
	"testing"
	"time"
)

func TestSaveUser(t *testing.T) {
	now := time.Now()
	serialized := false
	serializeUsers = func(users []User) error {
		serialized = true
		return nil
	}
	user := User{
		Username: "bob",
		Password: "hope",
	}
	user.Save()

	// Now check to ensure a created date was set on the user
	user = GetUser("bob")
	assert.True(t, user.Datecreated.After(now))
	assert.True(t, serialized)
}

func TestSaveUserWithSerializeFailure(t *testing.T) {
	serializeUsers = func(users []User) error {
		return errors.New("oh no!")
	}
	user := User{
		Username: "bob",
		Password: "hope",
	}
	err := user.Save()
	assert.NotNil(t, err)
}

func TestGetEmptyUser(t *testing.T) {
	user := GetUser("non-existing")
	assert.Nil(t, user.Token)
	assert.Equal(t, "", user.Username)
}

func TestGetUsers(t *testing.T) {
	users := GetUsers()
	assert.Equal(t, 1, len(users))
}

func TestDeleteEmptyUser(t *testing.T) {
	serialized := false
	serializeUsers = func(users []User) error {
		serialized = true
		return nil
	}
	err := DeleteUser("non-existing")
	assert.Nil(t, err)
	assert.True(t, serialized)
}
