package main

import (
	"github.com/sgoertzen/xchango"
	"google.golang.org/api/calendar/v3"
	"log"
	"time"
)

type EventActions struct {
	toAdd    []*calendar.Event
	toUpdate []*calendar.Event
	toDelete []*calendar.Event
}

func Sync(user User) {
	log.Printf("Syncing user %s", user.Username)
	user.State = syncing
	user.Save()

	appointments, err := xchango.GetExchangeAppointments(user.ExUser, user.ExCal)
	events, err := getGCalAppointments(user)
	if err == nil {
		err = mergeEvents(user, appointments, events)
		user.State = successfulsync
	}
	if err != nil {
		log.Println("Error while syncing events for", user, err)
		user.State = syncingerror
	}
	user.LastSync = time.Now()
	user.Save()
}

// TODO< make the appointment array a pointer as well
func mergeEvents(user User, appointments []xchango.Appointment, events *calendar.Events) error {

	actions, err := buildDiffLists(appointments, events)

	client := getClient(user)
	srv, err := calendar.New(client)
	if err != nil {
		log.Println("Unable to retrieve calendar Client %v", err)
		return err
	}

	for _, add := range actions.toAdd {
		log.Printf("Adding event of %s on %s", add.Summary, add.Start)
		_, err := srv.Events.Insert(user.GCalid, add).Do()
		if err != nil {
			log.Println(err)
			return err
		}
	}
	for _, edit := range actions.toUpdate {
		log.Printf("Updating event of %s on %s", edit.Summary, edit.Start)
		_, err := srv.Events.Update(user.GCalid, edit.Id, edit).Do()
		if err != nil {
			log.Println(err)
			return err
		}
	}
	for _, del := range actions.toDelete {
		log.Printf("Deleting event %s on %s", del.Summary, del.Start)
		err := srv.Events.Delete(user.GCalid, del.Id).Do()
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func buildDiffLists(appointments []xchango.Appointment, events *calendar.Events) (EventActions, error) {
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
	return eventActions, nil
}

func populateEvent(e *calendar.Event, a *xchango.Appointment) bool {
	var changes = false

	if e.Summary != a.Subject {
		e.Summary = a.Subject
		changes = true
		log.Println("Summary has changed: ", e.Summary, a.Subject)
	}

	if e.Location != a.Location {
		e.Location = a.Location
		changes = true
		log.Println("Location has changed: ", e.Location, a.Location)
	}

	desc := a.BuildDesc()
	if e.Description != desc {
		e.Description = desc
		changes = true
		log.Println("Description has changed")
	}

	if a.IsAllDayEvent {
		eventStart := calendar.EventDateTime{Date: a.Start.Format("2006-01-02")}
		if e.Start == nil || e.Start.Date != eventStart.Date {
			log.Println("Start has changed: ", eventStart.Date)
			e.Start = &eventStart
			e.End = &eventStart
			changes = true
		}
	} else {
		appStart := a.Start.UTC().Format(time.RFC3339)
		appEnd := a.End.UTC().Format(time.RFC3339)

		var eventStart, eventEnd string
		if e.Start != nil {
			timed, _ := time.Parse(time.RFC3339, e.Start.DateTime)
			eventStart = timed.UTC().Format(time.RFC3339)
		}
		if appStart != eventStart {
			log.Println("Start has changed: ", eventStart, appStart)
			e.Start = &calendar.EventDateTime{DateTime: appStart}
			changes = true
		}

		if e.End != nil {
			timed, _ := time.Parse(time.RFC3339, e.Start.DateTime)
			eventEnd = timed.UTC().Format(time.RFC3339)
		}
		if appEnd != eventEnd {
			e.End = &calendar.EventDateTime{DateTime: appEnd}
			// Don't record change here as google adjusts the end time between 5 and 15 minutes less.
		}
	}
	if changes {
		log.Printf("Changes found on appointment %s", a.Subject)
	}

	e.ExtendedProperties = &calendar.EventExtendedProperties{
		Private: map[string]string{"ItemId": a.ItemId},
	}
	return changes
}
