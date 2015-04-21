package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCanary(t *testing.T) {
	assert.True(t, true, "This is good. Canary test passing")
}

func TestSimpleEncryption(t *testing.T) {
	testEncryption(t, "a very very very very secret key", "String to encrypt")
}

func TestSimpleDecryption(t *testing.T) {
	testDecryption(t, "a very very very very secret key", "9v9CTyYgqz7OlqwkNFcO4nEMAI7tTH5gMCuGYaH6VLNj", "String to encrypt")
}

func TestKey16(t *testing.T) {
	testEncryption(t, "1234567890123456", "String to encrypt")
}

func testDecryption(t *testing.T, key string, ciphertext string, plaintext string) {
	keybytes := []byte(key)
	cipherbytes := []byte(ciphertext)
	decryptedtext, err := Decrypt(keybytes, cipherbytes)
	assert.Nil(t, err)
	assert.NotNil(t, decryptedtext)
	assert.Equal(t, decryptedtext, plaintext)
}

func testEncryption(t *testing.T, key string, plaintext string) {
	keybytes := []byte(key)
	plainbytes := []byte(plaintext)
	encrypted, err := Encrypt(keybytes, plainbytes)
	assert.Nil(t, err)
	assert.NotNil(t, encrypted)
}
