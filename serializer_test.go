package main

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestSetKeyShort(t *testing.T) {
  err := SetKey("short")
  assert.NotNil(t, err)
}

func TestSetKey16(t *testing.T) {
  err := SetKey("_sixteen chars._")
  assert.Nil(t, err)
}

func TestSetKey20(t *testing.T) {
  err := SetKey("twenty chars14161820")
  assert.Nil(t, err)
  key := Key()
  assert.Equal(t, string(key), "twenty chars14161820****")
}
