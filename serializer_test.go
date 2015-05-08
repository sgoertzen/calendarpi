package main

import (
	"github.com/sgoertzen/xchango"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializeAndDeserialize(t *testing.T) {
	users := []User{
		User{
			ExUser: &xchango.ExchangeUser{
				Username: "green",
				Password: "red",
			},
			ExCal: &xchango.ExchangeCalendar{
				Folderid: "fid",
				Changekey: "ckey",
			},
		},
	}
	key, _ := CreateKey("samplekey19385558")
	setStoredKey(key)
	serialized, err := serializeAndEncrypt(users)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	users2, err := decryptAndDeserialize(key, []byte(serialized))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users2))
	assert.Equal(t, "green", users2[0].ExUser.Username)
	assert.Equal(t, "red", users2[0].ExUser.Password)
	assert.Equal(t, "fid", users2[0].ExCal.Folderid)
	assert.Equal(t, "ckey", users2[0].ExCal.Changekey)
}
