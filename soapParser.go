package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"time"
)

type Organizer struct {
	Mailbox Mailbox
}

type Mailbox struct {
	Name string
}

type ItemId struct {
	Id        string `xml:"Id,attr"`
	ChangeKey string `xml:"ChangeKey,attr"`
}

type CalendarItem struct {
	ItemId         ItemId
	Subject        string
	DisplayCc      string
	DisplayTo      string
	Start          string
	End            string
	IsAllDayEvent  bool
	Location       string
	MyResponseType string
	Organizer      Organizer
}

type Appointment struct {
	ItemId         string
	Subject        string
	Cc             string
	To             string
	Start          time.Time
	End            time.Time
	IsAllDayEvent  bool
	Location       string
	MyResponseType string
	Organizer      string
}

// TODO: refactor this method with ParseAppointments
func ParseCalendarFolder(soap string) ItemId {
	// TODO: Should I just pass in a byte[] instead of string
	decoder := xml.NewDecoder(bytes.NewBufferString(soap))

	//itemId := make(ItemId)
	var itemId ItemId

	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "FolderId" {
				decoder.DecodeElement(&itemId, &se)
				break
				//appointments = append(appointments, item.ToAppointment())
			}
		}
	}
	return itemId
}

func ParseAppointments(soap string) []Appointment {
	// TODO: Should I just pass in a byte[] instead of string
	decoder := xml.NewDecoder(bytes.NewBufferString(soap))

	appointments := make([]Appointment, 0)

	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "CalendarItem" {
				var item CalendarItem
				decoder.DecodeElement(&item, &se)
				appointments = append(appointments, item.ToAppointment())
			}
		}
	}
	return appointments
}

func (c CalendarItem) ToAppointment() Appointment {
	app := Appointment{
		ItemId:         c.ItemId.Id,
		Subject:        c.Subject,
		Cc:             c.DisplayCc,
		To:             c.DisplayTo,
		IsAllDayEvent:  c.IsAllDayEvent,
		Location:       c.Location,
		MyResponseType: c.MyResponseType,
		Organizer:      c.Organizer.Mailbox.Name,
	}
	t1, err := time.Parse(time.RFC3339, c.Start)
	if err != nil {
		log.Printf("Error while parsing time.  Time string is: ", c.Start, err)
	}
	app.Start = t1

	t1, err = time.Parse(time.RFC3339, c.End)
	if err != nil {
		log.Printf("Error while parsing time.  Time string is: ", c.End, err)
	}
	app.End = t1
	return app
}

func (a Appointment) String() string {
	return fmt.Sprintf("%s starting %d", a.Subject, a.Start)
}
