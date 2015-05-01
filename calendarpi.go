package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

/*
TODO: Make private key for user and embed as hidden on page.  Verify this before saving.
TODO: support adding a new calendar (To add a new calendar to srv.Insert(calendar) - returns CalendarsInsertCall)
*/

const configfile = "conf.json"

func main() {
	config := readConfig()
	SetOauthConfig(config)
	SetExchangeConfig(config)
	go RunServer(config)
	runSyncLoop(config)
	log.Println("Server is exiting")
}

func readConfig() Config {

	file, err := os.Open("conf.json")
	if err != nil {
		if os.IsExist(err) {
			log.Fatal(err)
		}
		data, _ := ioutil.ReadFile("conf.template.json")
		ioutil.WriteFile("conf.json", data, 777)
		log.Fatal("No configuration file found.  A new configuration file has automatically been created for you.  Please edit conf.json and fill in the correct values.")
	}
	decoder := json.NewDecoder(file)
	config := Config{}
	if err = decoder.Decode(&config); err != nil {
		log.Fatal("Unable to parse the configuration file 'conf.json'", err)
	}
	return config
}

func runSyncLoop(config Config) {
	for true {
		log.Printf("Sleeping for %d minute(s)", config.MinutesBetweenSync())
		sleepTime := time.Duration(config.MinutesBetweenSync() * 60 * 1e9)
		time.Sleep(sleepTime)
		log.Println("Waking")
		users := GetUsers()
		for _, user := range users {
			Sync(user)
		}
	}
}
