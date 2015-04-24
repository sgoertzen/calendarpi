package main

import (
	"google.golang.org/api/calendar/v3"
	"log"
	"time"
)

func GetCalendarList(user User) *calendar.CalendarList {
	client := getClient(user)
	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to get calendar service", err)
	}
	calendars, err2 := srv.CalendarList.List().Do()
	if err2 != nil {
		log.Fatalf("Unable to get calendar list", err)
	}
	return calendars

	// To add a new calendar to srv.Insert(calendar) - returns CalendarsInsertCall
}

func processAppointments(user User, apps []Appointment) {
	//for _, app := range apps {
	//log.Println(app.Subject)

	//}
}

// TODO, max results should be part of the app config
// TODO, move into another class
func getGCalAppointments(user User, gcalId string) error {
	client := getClient(user)
	srv, err := calendar.New(client)
	if err != nil {
		// TODO: switch failures to Fatals
		log.Fatalf("Unable to retrieve calendar Client %v", err)
	}
	log.Println("Got client: ", client)
	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).
		TimeMin(t).MaxResults(100).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve the user's events. %v", err)
	}
	log.Println(events)
	log.Println(events.Items)
	return nil
}
func mergeEvents(appointments []Appointment) {
	// Event object is described here: https://godoc.org/google.golang.org/api/calendar/v3#Event
	/*for _, i := range events.Items {
		var when string
		// If the DateTime is an empty string the Event is an all-day Event.
		// So only Date is available.
		if i.Start.DateTime != "" {
			when = i.Start.DateTime
		} else {
			when = i.Start.Date
		}
		fmt.Printf("%s (%s)\n", i.Summary, when)
		// TODO Look at these
		i.Description
		i.End.DateTime
		i.End.Date
		i.Etag
		i.ExtendedProperties

	}*/
}
