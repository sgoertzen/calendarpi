package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
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
		return errors.New("Unable to serialize")
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

func TestGetUsersSorted(t *testing.T) {
	userZ := User{
		Username: "zzz",
		Password: "zzz",
	}
	userA := User{
		Username: "aaa",
		Password: "aaa",
	}
	userZ.Save()
	userA.Save()
	users := GetUsers()
	assert.Equal(t, 3, len(users))
	assert.Equal(t, "aaa", users[0].Username)
	assert.Equal(t, "bob", users[1].Username)
	assert.Equal(t, "zzz", users[2].Username)
}
