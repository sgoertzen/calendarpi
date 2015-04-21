package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCanary(t *testing.T) {
	assert.True(t, true, "This is good. Canary test passing")
}

func TestCreateKeyShort(t *testing.T) {
	key, err := CreateKey("short")
	assert.NotNil(t, err)
	assert.Equal(t, []byte(nil), key)
}

func TestCreateKey16(t *testing.T) {
	key, err := CreateKey("_sixteen chars._")
	assert.Nil(t, err)
	assert.Equal(t, "_sixteen chars._", string(key))
}

// This should create a key that is 32 long
func TestCreateKey18(t *testing.T) {
	key, err := CreateKey("twenty chars141618")
	assert.Nil(t, err)
	assert.Equal(t, "twenty chars141618**************", string(key))
}

// This should truncate to 32 characters
func TestCreateKeyLarge(t *testing.T) {
	key, err := CreateKey("akfljeaiawvelaewjaewmlkawklmaewelkaewljjjkhaewkjvenk")
	assert.Nil(t, err)
	assert.Equal(t, "akfljeaiawvelaewjaewmlkawklmaewe", string(key))
}

func TestSimpleDecryption(t *testing.T) {
	testDecryption(t, "a very very very very secret key", "9v9CTyYgqz7OlqwkNFcO4nEMAI7tTH5gMCuGYaH6VLNj", "String to encrypt")
}

func testDecryption(t *testing.T, key string, ciphertext string, plaintext string) {
	keybytes := []byte(key)
	cipherbytes := []byte(ciphertext)
	decryptedtext, err := Decrypt(keybytes, cipherbytes)
	assert.Nil(t, err)
	assert.NotNil(t, decryptedtext)
	assert.Equal(t, decryptedtext, plaintext)
}

func TestRoundTrip(t *testing.T) {
	var tests = []struct {
		key       string
		plaintext string
	}{
		{"a very very very very secret key", "1. String to encrypt"},
		{"a very very secret key", "1. String to encrypt"},
		{"1234567890123456", "2. String to encrypt"},
		//{"pvR@(@cnkj332*@(@ECNELJsevi9ryu$A2keX%w3qCEq28lbqnDPpW21q937cniu", "3. String to encrypt"},
	}
	for _, test := range tests {
		keybytes, err := CreateKey(test.key)
		assert.Nil(t, err)
		plainbytes := []byte(test.plaintext)
		t.Log(len(keybytes))
		encrypted, err := Encrypt(keybytes, plainbytes)
		assert.Nil(t, err)
		assert.NotNil(t, encrypted)

		decrypted, err := Decrypt(keybytes, []byte(encrypted))
		assert.Equal(t, test.plaintext, decrypted)
	}

}
