package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
  err := DeleteUser("non-existing")
	assert.Nil(t, err)
}
