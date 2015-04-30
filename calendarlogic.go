package main

import (
	"github.com/kennygrant/sanitize"
	"google.golang.org/api/calendar/v3"
	"log"
	"time"
)

type EventActions struct {
	toAdd    []*calendar.Event
	toUpdate []*calendar.Event
	toDelete []*calendar.Event
}

func mergeEvents(user User, appointments []Appointment, events *calendar.Events) error {

	actions, err := buildDiffLists(appointments, events)

	client := getClient(user)
	srv, err := calendar.New(client)
	if err != nil {
		log.Println("Unable to retrieve calendar Client %v", err)
		return err
	}

	for _, add := range actions.toAdd {
		log.Println("Adding event of ", add.Summary)
		_, err := srv.Events.Insert(user.GCalid, add).Do()
		if err != nil {
			log.Println(err)
			return err
		}
	}
	for _, edit := range actions.toUpdate {
		log.Println("Updating event of ", edit.Summary)
		_, err := srv.Events.Update(user.GCalid, edit.Id, edit).Do()
		if err != nil {
			log.Println(err)
			return err
		}
	}
	for _, del := range actions.toDelete {
		log.Println("Deleting event " + del.Summary)
		err := srv.Events.Delete(user.GCalid, del.Id).Do()
		if err != nil {
			log.Println(err)
			return err
		}
	}
	user.LastSync = time.Now()
	user.Save()
	return nil
}

func buildDiffLists(appointments []Appointment, events *calendar.Events) (EventActions, error) {
	var itemMap = make(map[string]*calendar.Event)
	for _, event := range events.Items {
		if event.ExtendedProperties == nil || len(event.ExtendedProperties.Private["ItemId"]) == 0 {
			continue
		}
		itemMap[event.ExtendedProperties.Private["ItemId"]] = event
	}

	var eventActions EventActions
	for _, app := range appointments {
		existingEvent := itemMap[app.ItemId]
		if existingEvent != nil {
			delete(itemMap, app.ItemId)
			//newEvent := calendar.Event{}
			//newEvent.Id = existingEvent.Id
			changes := populateEvent(existingEvent, &app)
			if changes {
				eventActions.toUpdate = append(eventActions.toUpdate, existingEvent)
			}
			continue
		}
		e := calendar.Event{}
		populateEvent(&e, &app)
		eventActions.toAdd = append(eventActions.toAdd, &e)
	}
	for _, e := range itemMap {
		eventActions.toDelete = append(eventActions.toDelete, e)
	}
	log.Println("Total events: ", len(events.Items))
	log.Println("Count of updated events: ", len(eventActions.toUpdate))
	return eventActions, nil
}

func populateEvent(e *calendar.Event, a *Appointment) bool {
	var changes = false

	if e.Summary != a.Subject {
		log.Println("Subjects are different.  Summary vs Subject ", e.Summary, a.Subject)
		e.Summary = a.Subject
		changes = true
	}

	if e.Location != a.Location {
		log.Println("Locations are different.  GCal vs Exchange ", e.Location, a.Location)
		e.Location = a.Location
		changes = true
	}

	desc := buildDesc(a)
	if e.Description != desc {
		log.Println("Descriptions are different.  GCal vs Exchange ", e.Description, desc)
		e.Description = desc
		changes = true
	}

	var eventStart, eventEnd calendar.EventDateTime
	if a.IsAllDayEvent {
		eventStart = calendar.EventDateTime{Date: a.Start.Format("2006-01-02")}
		if e.Start != nil && e.Start.Date != eventStart.Date {
			log.Println("Starts are different.  GCal vs Exchange ")
			log.Println("GCal ", e.Start.Date)
			log.Println("Exchange ", eventStart.Date)
			changes = true
		}
		e.Start = &eventStart
	} else {
		eventStart = calendar.EventDateTime{DateTime: a.Start.Format(time.RFC3339)}
		eventEnd = calendar.EventDateTime{DateTime: a.End.Format(time.RFC3339)}
		if e.Start != nil && e.Start.DateTime != eventStart.DateTime {
			log.Println("Starts are different.  GCal vs Exchange ")
			log.Println("GCal ", e.Start.Date)
			log.Println("Exchange ", eventStart.Date)
			changes = true
		}
		e.Start = &eventStart
		e.End = &eventEnd

	}

	//e.End != &eventEnd ||
	e.ExtendedProperties = &calendar.EventExtendedProperties{
		Private: map[string]string{"ItemId": a.ItemId},
	}
	return changes
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
