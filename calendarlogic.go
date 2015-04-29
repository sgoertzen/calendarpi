package main

import (
	"github.com/kennygrant/sanitize"
	"google.golang.org/api/calendar/v3"
	"log"
	"time"
)

type EventActions struct {
	toAdd []*calendar.Event
	toUpdate []*calendar.Event
	toDelete []*calendar.Event
}

func mergeEvents(user User, appointments []Appointment, events *calendar.Events) error {

	actions, err := buildDiffLists(appointments, events)

	client := getClient(user)
	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
		return err
	}

	for _, event := range actions.toAdd {
		_, err := srv.Events.Insert(user.GCalid, event).Do()
		if err != nil {
			log.Println(err)
			return err
		}
	}
	for _, edit := range actions.toUpdate {
		_, err := srv.Events.Patch(user.GCalid, edit.Id, edit).Do()
		if err != nil {
			log.Println(err)
			return err
		}
	}
	for _, del := range actions.toDelete {
		err := srv.Events.Delete(user.GCalid, del.Id).Do()
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func buildDiffLists(appointments []Appointment, events *calendar.Events) (EventActions, error) {
	var itemMap = make(map[string]*calendar.Event)
	for _, event := range events.Items {
		if event.ExtendedProperties == nil || len(event.ExtendedProperties.Private["ItemId"]) == 0{
			// Skip this as it isn't one of our calendar appointments
			continue
		}
		itemMap[event.ExtendedProperties.Private["ItemId"]] = event
	}

	var eventActions EventActions
	for _, app := range appointments {
		existingEvent := itemMap[app.ItemId]
		if existingEvent != nil {
			log.Println("Skipping due to appointment already existing")
			// todo remove from map
			delete(itemMap, app.ItemId)
			eventActions.toUpdate = append(eventActions.toUpdate, existingEvent)
			continue
		}
		e := calendar.Event{}
		populateEvent(&e, &app)
		eventActions.toAdd = append(eventActions.toAdd, &e)
	}
	for _, e := range itemMap {
		// todo put into del map
		eventActions.toDelete = append(eventActions.toDelete, e)
	}
	return eventActions, nil
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
