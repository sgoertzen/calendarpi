package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var backupFile = "blabbersnatzle.bak"

// TODO: Read this in on first start and store in memory
var key = []byte("example key 1234")


func Save() error {
	log.Println("Serializing users")
	os.Remove(backupFile)
	users := GetUsers()
	data, err := json.Marshal(users)
	if err != nil {
		log.Println("Unable to json the users!")
		return err
	}
	encryptedData := Encrypt(key, data)
	err2 := ioutil.WriteFile(backupFile, []byte(encryptedData), 777)
	if err2 != nil {
		log.Println("Unable to backup the file!")
		return err2
	}
	return nil
}

func Load() error {
	log.Println("Unserializing users")
	filebytes, err := ioutil.ReadFile(backupFile)
	if err != nil {
		log.Println("Unable to find the backed up users file.")
		return err
	}
	decryptedData := Decrypt(key, string(filebytes))
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
