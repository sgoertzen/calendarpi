package main

import (
	"crypto/aes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var backupFile = "blabbersnatzle.bak"

var key []byte

func Key() []byte {
	return key
}

func SetKey(keystring string) error {
	if len(keystring) < aes.BlockSize {
		message := fmt.Sprintf("Key is to short!  Must be at least %d", aes.BlockSize)
		return errors.New(message)
	}
	if len(keystring)%aes.BlockSize != 0 {
		keystring = keystring + strings.Repeat("*", len(keystring)%aes.BlockSize)
	}
	key = []byte(keystring)
	return nil
}

func SerializeUsers() error {
	log.Println("Serializing users")
	os.Remove(backupFile)
	users := GetUsers()
	data, err := json.Marshal(users)
	if err != nil {
		log.Println("Unable to json the users!")
		return err
	}
	encryptedData, err := Encrypt(key, data)
	if err != nil {
		log.Println("Unable to encyrpt the data")
		return err
	}
	err2 := ioutil.WriteFile(backupFile, []byte(encryptedData), 777)
	if err2 != nil {
		log.Println("Unable to backup the file!")
		return err2
	}
	return nil
}

func DeserializeUsers() error {
	log.Println("Unserializing users")
	filebytes, err := ioutil.ReadFile(backupFile)
	if err != nil {
		log.Println("Unable to find the backed up users file.")
		return err
	}
	decryptedData, err := Decrypt(key, filebytes)
	if err != nil {
		log.Println("Unable to decrypt the file")
		return err
	}
	var users []User
	err3 := json.Unmarshal([]byte(decryptedData), &users)
	if err3 != nil {
		log.Println("Error while unmarshaling the json")
		return err
	}
	err2 := os.Remove(backupFile)
	if err2 != nil {
		log.Println("Unable to delete the file!")
		return err
	}
	for _, user := range users {
		user.Save()
	}
	return nil
}
