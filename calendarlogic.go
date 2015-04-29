package main

import (
	"github.com/kennygrant/sanitize"
	"google.golang.org/api/calendar/v3"
	"log"
	"time"
)

func mergeEvents(user User, appointments []Appointment, events *calendar.Events) error {

	addEvents, editEvents, err := buildDiffLists(appointments, events)

	client := getClient(user)
	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
		return err
	}
	log.Println(srv)
	for _, event := range addEvents {
		retevent, err := srv.Events.Insert(user.GCalid, event).Do()
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println(retevent)
	}
	for _, edit := range editEvents {
		//log.Println("Would update event ", edit.Summary)
		srv.Events.Patch(user.GCalid, edit.Id, edit).Do()
	}
	return nil
}

func buildDiffLists(appointments []Appointment, events *calendar.Events) ([]*calendar.Event, []*calendar.Event, error) {

	var itemMap = make(map[string]*calendar.Event)
	for _, event := range events.Items {
		itemMap[event.ExtendedProperties.Private["ItemId"]] = event
	}

	var addEvents []*calendar.Event
	var editEvents []*calendar.Event
	log.Printf("Looping over %d appointments", len(appointments))
	for _, app := range appointments {
		existingEvent := itemMap[app.ItemId]
		if existingEvent != nil {
			log.Println("Skipping due to appointment already existing")
			editEvents = append(editEvents, existingEvent)
			continue
		}
		event := calendar.Event{}
		populateEvent(&event, &app)
		addEvents = append(addEvents, &event)
	}
	return addEvents, editEvents, nil
}

func populateEvent(e *calendar.Event, a *Appointment) {
	var eventStart, eventEnd calendar.EventDateTime
	if a.IsAllDayEvent {
		eventStart = calendar.EventDateTime{Date: a.Start.Format("2006-01-02")}
	} else {
		eventStart = calendar.EventDateTime{DateTime: a.Start.Format(time.RFC3339)}
		eventEnd = calendar.EventDateTime{DateTime: a.End.Format(time.RFC3339)}
	}

	e.Summary = a.Subject
	e.Location = a.Location
	e.Start = &eventStart
	e.End = &eventEnd
	e.Description = buildDesc(a)
	e.ExtendedProperties = &calendar.EventExtendedProperties{
		Private: map[string]string{"ItemId": a.ItemId},
	}
}

func buildDesc(a *Appointment) string {
	var desc = ""

	addField := func(field string, label string) {
		if len(field) > 0 {
			desc += label + " " + field + "\n"
		}
	}
	addField(a.Organizer, "Organizer:")
	addField(a.To, "To:")
	addField(a.Cc, "Cc:")
	addField(a.MyResponseType, "Response:")
	body := sanitize.HTML(a.Body)
	desc += "\n" + body
	return desc
}
