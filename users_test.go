package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


func mock_Serializer(users []User) error {
	// do I need to do anything here?
	return nil
}

func TestGetEmptyUser(t *testing.T) {
	user := GetUser("non-existing")
	assert.Nil(t, user.Token)
	assert.Equal(t, "", user.Username)
}

func TestGetEmptyUsers(t *testing.T) {
	users := GetUsers()
	assert.Equal(t, 0, len(users))
}

func TestDeleteEmptyUser(t *testing.T) {
	MySerializeUsers = mock_Serializer
	err := DeleteUser("non-existing")
	assert.Nil(t, err)
}
