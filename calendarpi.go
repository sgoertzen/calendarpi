package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const configfile = "conf.json"

func main() {
	config := readConfig()
	
	f, err := os.OpenFile("output.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
	    log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()  // This line must be done in this method
	log.SetOutput(f)
	
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
	log.Printf("Sync Interval: %s", config.TimeBetweenSync())
	syncInterval, err := time.ParseDuration(config.TimeBetweenSync())
	if err != nil {
		log.Fatal("Unable to parse sleep time", err)
	}
		
	sleepTime, err := time.ParseDuration("1m")
	for true {	
		log.Println("syncing")
		users := GetUsers()
		for _, user := range users {
			if time.Now().After(user.LastSync.Add(syncInterval)) {
				Sync(user)
			}
		}
		time.Sleep(sleepTime)
	}
}
