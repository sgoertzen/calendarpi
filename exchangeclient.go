package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

type ExchangeConfig interface {
	ExchangeURL() string
}

var exchangeURL string

func SetConfig2(config ExchangeConfig) {
	exchangeURL = config.ExchangeURL()
}

func getFolderAndChangeKey(user User) User {

	soapReq, err := ioutil.ReadFile("xml/soapFolderRequest.xml")
	if err != nil {
		log.Printf("Error is", err)
	}
	results, err := postContents(soapReq, user)

	if err != nil {
		log.Println(err)
	}

	item := ParseCalendarFolder(string(results))

	user.Folderid = item.Id
	user.Changekey = item.ChangeKey
	user.State = exchangeLoginVerified
	log.Println("Fetched folderid of ", user.Folderid)
	user.Save()
	return user
}

func getExchangeCalendarData(user User) string {
	soapRequest := buildCalendarRequest(user.Folderid, user.Changekey)
	soapResponse, err := postContents(soapRequest, user)
	if err != nil {
		log.Println("Error while getting soap response")
		log.Fatal(err)
	}
	return soapResponse
}

func buildCalendarRequest(folderid string, changekey string) []byte {
	// TODO: Make dates use current date and two weeks in the future
	startDate := "2015-02-20T17:30:24.127Z"
	endDate := "2015-04-20T17:30:24.127Z"

	data := struct {
		StartDate string
		EndDate   string
		FolderId  string
		ChangeKey string
	}{
		startDate,
		endDate,
		folderid,
		changekey,
	}

	t, err := template.ParseFiles("xml/soapCalendarRequest.xml")
	if err != nil {
		log.Printf("Error is", err)
	}
	var doc bytes.Buffer
	t.Execute(&doc, data)
	if err != nil {
		log.Printf("Error while building contents ", err)
	}

	return doc.Bytes()
}

func postContents(contents []byte, user User) (string, error) {
	req2, err := http.NewRequest("POST", exchangeURL, bytes.NewBuffer(contents))
	req2.Header.Set("Host", user.Username+"@webmail.vwgoa.com")
	req2.Header.Set("Content-Type", "text/xml")
	req2.SetBasicAuth("na/"+user.Username, user.Password)

	client := &http.Client{}
	response, err := client.Do(req2)
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	return string(content), nil
}
