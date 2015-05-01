package main

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/calendar/v3"
	"testing"
	"time"
)

func TestBuildDiffListsEmpty(t *testing.T) {
	apps := []Appointment{}
	events := calendar.Events{}

	actions, err := buildDiffLists(apps, &events)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(actions.toAdd))
	assert.Equal(t, 0, len(actions.toUpdate))
}

func TestBuildDiffListsAllNew(t *testing.T) {
	apps := []Appointment{
		Appointment{
			ItemId:   "Hello",
			Subject:  "sub",
			Location: "loc",
			Body:     "<html><body><b>body</b></body></html>",
		},
	}
	events := calendar.Events{}

	actions, err := buildDiffLists(apps, &events)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(actions.toAdd))
	assert.Equal(t, "sub", actions.toAdd[0].Summary)
	assert.Equal(t, "loc", actions.toAdd[0].Location)
	assert.Equal(t, "\nbody", actions.toAdd[0].Description)
	assert.Equal(t, 0, len(actions.toUpdate))
	assert.Equal(t, 0, len(actions.toDelete))
}

func TestBuildDiffListsDelete(t *testing.T) {
	apps := []Appointment{}
	events := calendar.Events{
		Items: []*calendar.Event{
			&calendar.Event{
				Id: "45",
				ExtendedProperties: &calendar.EventExtendedProperties{
					Private: map[string]string{"ItemId": "45"},
				},
			},
		},
	}

	actions, err := buildDiffLists(apps, &events)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(actions.toDelete))
	assert.Equal(t, "45", actions.toDelete[0].Id)
	assert.Equal(t, 0, len(actions.toUpdate))
	assert.Equal(t, 0, len(actions.toAdd))
}

func TestBuildDiffListsLeaveExisting(t *testing.T) {
	apps := []Appointment{}
	events := calendar.Events{
		Items: []*calendar.Event{
			&calendar.Event{
				Id: "45",
				// Note: No extended properties on here.  Indicates not our event.
			},
		},
	}

	actions, err := buildDiffLists(apps, &events)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(actions.toDelete))
}

func TestPopulateEventEmpty(t *testing.T) {
	e := calendar.Event{}
	a := Appointment{
		Subject: "42",
		ItemId:  "uniqueId",
	}
	populateEvent(&e, &a)
	assert.Equal(t, "42", e.Summary)
	assert.Equal(t, "uniqueId", e.ExtendedProperties.Private["ItemId"])
}

func TestPopulateEventExisting(t *testing.T) {
	e := calendar.Event{
		Id:      "train",
		Summary: "blah",
	}
	a := Appointment{
		Subject: "42",
		ItemId:  "uniqueId",
	}
	changes := populateEvent(&e, &a)
	assert.Equal(t, true, changes)
	assert.Equal(t, "42", e.Summary)
	assert.Equal(t, "train", e.Id)
	assert.Equal(t, "uniqueId", e.ExtendedProperties.Private["ItemId"])
}

func TestPopulateEventExistingNoChanges(t *testing.T) {
	e := calendar.Event{
		Id:          "123",
		Summary:     "phone call",
		Description: "\nhello",
		Start:       &calendar.EventDateTime{DateTime: "2015-05-04T11:00:00-07:00"},
		End:         &calendar.EventDateTime{DateTime: "2015-05-04T12:00:00-07:00"},
	}
	tStart, _ := time.Parse(time.RFC3339, "2015-05-04T18:00:00Z")
	tEnd, _ := time.Parse(time.RFC3339, "2015-05-04T18:00:00Z")
	a := Appointment{
		Subject:       "phone call",
		ItemId:        "123",
		Start:         tStart,
		End:           tEnd,
		IsAllDayEvent: false,
		Body:          "hello",
	}
	changes := populateEvent(&e, &a)
	assert.False(t, changes)
}

func TestPopulateEventExistingStartChange(t *testing.T) {
	e := calendar.Event{
		Id:          "123",
		Summary:     "phone call",
		Description: "\n",
		Start:       &calendar.EventDateTime{DateTime: "2015-04-12T16:00:00Z"},
	}
	t1, _ := time.Parse(time.RFC3339, "2015-04-13T16:00:00Z")
	a := Appointment{
		Subject:       "phone call",
		ItemId:        "123",
		Start:         t1,
		IsAllDayEvent: false,
	}
	changes := populateEvent(&e, &a)
	assert.True(t, changes)
}

func TestPopulateEventExistingEndChange(t *testing.T) {
	e := calendar.Event{
		Id:          "123",
		Summary:     "phone call",
		Description: "\n",
		Start:       &calendar.EventDateTime{DateTime: "2015-04-13T16:00:00Z"},
		End:         &calendar.EventDateTime{DateTime: "2015-04-13T18:00:00Z"},
	}
	tStart, _ := time.Parse(time.RFC3339, "2015-04-13T16:00:00Z")
	tEnd, _ := time.Parse(time.RFC3339, "2015-04-13T17:00:00Z")
	a := Appointment{
		Subject:       "phone call",
		ItemId:        "123",
		Start:         tStart,
		End:           tEnd,
		IsAllDayEvent: false,
	}
	changes := populateEvent(&e, &a)

	// Google adjusts the end time between 5 and 15 minutes.  We need to ignore end time differences.
	assert.False(t, changes)
}