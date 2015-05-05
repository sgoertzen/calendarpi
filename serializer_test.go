package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializeAndDeserialize(t *testing.T) {
	users := []User {
		User { Username: "green", },
	}
	key, _ := CreateKey("samplekey19385558")
	setStoredKey(key)
	serialized, err := serializeAndEncrypt(users)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)
	
	users2, err := decryptAndDeserialize(key, []byte(serialized))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users2))
	assert.Equal(t, "green", users2[0].Username)
}