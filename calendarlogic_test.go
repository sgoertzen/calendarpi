package main

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/calendar/v3"
	"testing"
)

func TestBuildDiffListsEmpty(t *testing.T) {
	apps := []Appointment{}
	events := calendar.Events{}

	addEvents, editEvents, err := buildDiffLists(apps, &events)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(addEvents))
	assert.Equal(t, 0, len(editEvents))
}

func TestBuildDiffListsAllNew(t *testing.T) {
	apps := []Appointment{
		Appointment{
			ItemId:   "Hello",
			Subject:  "sub",
			Location: "loc",
      Body: "<html><body><b>body</b></body></html>",
		},
	}
	events := calendar.Events{}

	addEvents, editEvents, err := buildDiffLists(apps, &events)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(addEvents))
	assert.Equal(t, "sub", addEvents[0].Summary)
	assert.Equal(t, "loc", addEvents[0].Location)
	assert.Equal(t, "\nbody", addEvents[0].Description)
	assert.Equal(t, 0, len(editEvents))
}

func TestPopulateEventEmpty(t *testing.T) {
	e := calendar.Event {}
	a := Appointment{
			Subject:  "42",
			ItemId: "uniqueId",
		}
	populateEvent(&e, &a)
	assert.Equal(t, "42", e.Summary)
	assert.Equal(t, "uniqueId", e.ExtendedProperties.Private["ItemId"])
}

func TestPopulateEventExisting(t *testing.T) {
	e := calendar.Event {
		Id: "train",
		Summary: "blah",
	}
	a := Appointment{
		Subject:  "42",
		ItemId: "uniqueId",
	}
	populateEvent(&e, &a)
	assert.Equal(t, "42", e.Summary)
	assert.Equal(t, "train", e.Id)
	assert.Equal(t, "uniqueId", e.ExtendedProperties.Private["ItemId"])
}

func TestBuildBody(t *testing.T) {
expected := 	`Organizer: minifig
To: deadpool
Cc: batman

body
body2`

	app := Appointment{
			To:  "deadpool",
			Cc:  "batman",
			Organizer: "minifig",
			Body: "<html><body><b>body</b><br/>body2</body></html>",
		}
	desc := buildDesc(&app)
	assert.Equal(t, expected, desc)
}
