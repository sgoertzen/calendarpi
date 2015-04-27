package main

import (
	"google.golang.org/api/calendar/v3"
	"log"
	"time"
)

func GetGCalendarList(user User) *calendar.CalendarList {
	client := getClient(user)
	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to get calendar service", err)
	}
	calendars, err2 := srv.CalendarList.List().MinAccessRole("writer").Do()
	if err2 != nil {
		log.Fatalf("Unable to get calendar list", err)
	}
	return calendars
}

// TODO, max results should be part of the app config
// TODO, move into another class
func getGCalAppointments(user User) (*calendar.Events, error) {
	client := getClient(user)
	srv, err := calendar.New(client)
	if err != nil {
		// TODO: switch failures to Fatals
		log.Fatalf("Unable to retrieve calendar Client %v", err)
	}
	log.Println("Got client: ", client)
	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List(user.GCalid).ShowDeleted(false).
		SingleEvents(true).
		TimeMin(t).MaxResults(100).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve the user's events. %v", err)
	}
	return events, nil
}

func mergeEvents(user User, appointments []Appointment, events *calendar.Events) error {
	client := getClient(user)
	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
		return err
	}
	for _, app := range appointments {
		eventExists := false
		for _, event := range events.Items {
			// TODO: THIS IS BAD.  Use a hashmap to search through the IDS.
			// Doing this quickly for now to see if anything actually works
			if event.ExtendedProperties.Private["ItemId"] == app.ItemId {
				log.Println("Found item of id: ", app.ItemId)
				eventExists = true
				break
			}
		}
		if eventExists {
			// TODO: Update the event in case it has changed some
			log.Println("Skipping due to appointment already existing")
			continue
		}

		var eventStart, eventEnd calendar.EventDateTime
		if app.IsAllDayEvent {
			eventStart = calendar.EventDateTime{
				Date: app.Start.Format("2006-01-02"),
			}
		} else {
			eventStart = calendar.EventDateTime{
				DateTime: app.Start.Format(time.RFC3339),
			}
			eventEnd = calendar.EventDateTime{
				DateTime: app.End.Format(time.RFC3339),
			}
		}

		event := calendar.Event{
			Summary:  app.Subject,
			Location: app.Location,
			Start:    &eventStart,
			//Id: app.ItemId,
			Description: app.ItemId,
			End:         &eventEnd,
			ExtendedProperties: &calendar.EventExtendedProperties{
				Private: map[string]string{"ItemId": app.ItemId},
			},
			//Transparency: appointment.MyResponseType == "Accepted"
		}
		//log.Println(event)
		retevent, err := srv.Events.Insert(user.GCalid, &event).Do()
		if err != nil {
			log.Println(err)
		}
		log.Println(retevent)
		return err
	}
	// Event object is described here: https://godoc.org/google.golang.org/api/calendar/v3#Event
	for _, i := range events.Items {
		var when string
		// If the DateTime is an empty string the Event is an all-day Event.
		// So only Date is available.
		if i.Start.DateTime != "" {
			when = i.Start.DateTime
		} else {
			when = i.Start.Date
		}
		log.Printf("%s (%s)\n", i.Summary, when)
		// TODO Look at these
		/*
			i.Description
			i.End.DateTime
			i.End.Date
			i.Etag
			i.ExtendedProperties
		*/
	}
	return nil
}
