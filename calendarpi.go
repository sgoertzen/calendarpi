package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

/*
TODO
Test output of soap request template
Test post call to server with soap request

Write job to actually do sync!
How are we going to show status of syncs and display errors
Allow sync to happen immediately on new accounts but scheduled after that

Give option to choose calendar or create a new one

Look into using mux for routing
Put files into different packages

Use this instead of err != nil all over
if _, err := io.ReadFull(rand.Reader, iv); err != nil {
	return "", err
}

*/

const configfile = "conf.json"

func main() {
	config := readConfig()
	SetConfig(config)
	SetConfig2(config)
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
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal("Unable to parse the configuration file 'conf.json'", err)
	}
	return config
}

func runSyncLoop(config Config) {
	for true {
		sleepTime := time.Duration(config.MinutesBetweenSync() * 60 * 1e9)
		time.Sleep(sleepTime)
		users := GetUsers()
		for _, user := range users {
			soapResults := getExchangeCalendarData(user)
			appointments := ParseAppointments(soapResults)
			log.Println("len:", len(appointments))
			processAppointments(user, appointments)
		}
	}
}
