package main

import (
	"google.golang.org/api/calendar/v3"
	"log"
	"time"
	"github.com/kennygrant/sanitize"
)

func mergeEvents(user User, appointments []Appointment, events *calendar.Events) error {

	addEvents, err := buildDiffLists(appointments, events)

	client := getClient(user)
	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
		return err
	}
	log.Println(srv)
	for _, event := range addEvents {
		retevent, err := srv.Events.Insert(user.GCalid, &event).Do()
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println(retevent)
	}
	return nil
}

func buildDiffLists(appointments []Appointment, events *calendar.Events) ([]calendar.Event, error) {

	var itemMap = make(map[string]*calendar.Event)
	for _, event := range events.Items {
		itemMap[event.ExtendedProperties.Private["ItemId"]] = event
	}

	var addEvents []calendar.Event
	log.Printf("Looping over %d appointments", len(appointments))
	for _, app := range appointments {
		existingEvent := itemMap[app.ItemId]
		if existingEvent != nil {
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
			Summary:     app.Subject,
			Location:    app.Location,
			Start:       &eventStart,
			//Description: StripTags(app.Body),
			// Could use this instead: https://github.com/kennygrant/sanitize
			Description: sanitize.HTML(app.Body),
			End: &eventEnd,
			ExtendedProperties: &calendar.EventExtendedProperties{
				Private: map[string]string{"ItemId": app.ItemId},
			},
		}
		addEvents = append(addEvents, event)
	}
	return addEvents, nil
}

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
	log.Printf("%s (%s)\n", i.Summary, when)
}*/
