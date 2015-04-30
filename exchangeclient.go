package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"time"
)

type ExchangeConfig interface {
	ExchangeURL() string
	MaxFetchSize() int
}

var exchangeConfig ExchangeConfig

func SetExchangeConfig(config ExchangeConfig) {
	exchangeConfig = config
}

func GetFolderAndChangeKey(user User) User {

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
	user.Save()
	return user
}

func GetExchangeAppointments(user User) []Appointment {
	// This first call will just get ids for each appt
	soapResults := getExchangeCalendarData(user)
	itemIds := ParseAppointments(soapResults)

	// This call will get all the fields given the ids
	soapResults = getExchangeAppointmentBodyData(itemIds, user)
	appointments := ParseAppointments(soapResults)

	return appointments
}

func getExchangeAppointmentBodyData(itemIds []Appointment, user User) string {
	soapRequest := buildCalendarDetailRequest(itemIds)
	soapResponse, err := postContents(soapRequest, user)
	if err != nil {
		log.Println("Error while getting soap response")
		log.Fatal(err)
	}
	return soapResponse
}

func buildCalendarDetailRequest(itemIds []Appointment) []byte {

	data := struct {
		Appointments []Appointment
	}{
		itemIds,
	}

	t, err := template.ParseFiles("xml/soapCalendarDetailRequest.xml")
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

func getExchangeCalendarData(user User) string {
	soapRequest := buildCalendarItemRequest(user.Folderid, user.Changekey)
	soapResponse, err := postContents(soapRequest, user)
	if err != nil {
		log.Println("Error while getting soap response", err)
		return ""
	}
	return soapResponse
}

func buildCalendarItemRequest(folderid string, changekey string) []byte {
	startDate := time.Now().UTC().Format(time.RFC3339)
	endDate := time.Now().UTC().AddDate(0, 0, 14).Format(time.RFC3339)

	data := struct {
		StartDate    string
		EndDate      string
		FolderId     string
		ChangeKey    string
		MaxFetchSize int
	}{
		startDate,
		endDate,
		folderid,
		changekey,
		exchangeConfig.MaxFetchSize(),
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
	req2, err := http.NewRequest("POST", exchangeConfig.ExchangeURL(), bytes.NewBuffer(contents))
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
