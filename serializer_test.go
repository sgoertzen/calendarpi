package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetKey(t *testing.T) {
	key := "1234567890123456"
	err := SetKey(key)
	assert.Nil(t, err)
	assert.Equal(t, key, string(Key()))
}
