package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

func CreateKey(keystring string) ([]byte, error) {
	if len(keystring) > 32 {
		keystring = keystring[:32]
	}

	// Key must be one of the following lengths: 16, 24, or 32
	if len(keystring) < aes.BlockSize {
		message := fmt.Sprintf("Key is to short!  Must be at least %d", aes.BlockSize)
		return nil, errors.New(message)
	}
	// Append to make an increment of the blocksize
	if len(keystring)%aes.BlockSize != 0 {
		keystring = keystring + strings.Repeat("*", aes.BlockSize-(len(keystring)%aes.BlockSize))
	}
	return []byte(keystring), nil
}

// encrypt string to base64 crypto using AES
func Encrypt(key []byte, plaintext []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// decrypt from base64 to decrypted string
func Decrypt(key []byte, cryptoText []byte) (string, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(string(cryptoText))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext), nil
}
