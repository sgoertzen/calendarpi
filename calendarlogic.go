package main

import (
	"google.golang.org/api/calendar/v3"
	"log"
	"time"
)

func mergeEvents(user User, appointments []Appointment, events *calendar.Events) error {
	client := getClient(user)
	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
		return err
	}
	log.Println("Looping over %s appointments", len(appointments))
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
		log.Println("Adding event named ", app.Subject)

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
			return err
		}
		log.Println(retevent)
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
