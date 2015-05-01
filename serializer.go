package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var backupFile = "blabbersnatzle.bak"

var storedKey []byte

func Key() []byte {
	return storedKey
}

func SerializeUsers(users []User) error {
	os.Remove(backupFile)
	encryptedData, err := serializeAndEncrypt(users)
	if err != nil {
		return err
	}
	err2 := ioutil.WriteFile(backupFile, []byte(encryptedData), 777)
	if err2 != nil {
		log.Println("Unable to backup the file!")
		return err2
	}
	return nil
}

func serializeAndEncrypt(users []User) (string, error) {
	data, err := json.Marshal(users)
	if err != nil {
		log.Println("Unable to json the users!")
		return "", err
	}
	encryptedData, err := Encrypt(storedKey, data)
	if err != nil {
		log.Println("Unable to encyrpt the data")
		return "", err
	}
	return encryptedData, nil
}

func DeserializeUsers(key string) error {
	localkey, err := CreateKey(key)
	log.Println("Unserializing users")
	filebytes, err := ioutil.ReadFile(backupFile)
	if err != nil {
		log.Println("Users file not present, skipping.") 
		storedKey = localkey
		return nil
	}
	decryptedData, err := Decrypt(localkey, filebytes)
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
	// Only save this key if we are successful
	storedKey = localkey
	
	for _, user := range users {
		user.Save()
	}
	return nil
}
