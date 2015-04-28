package main

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/calendar/v3"
	"testing"
)

func TestBuildDiffListsEmpty(t *testing.T) {
	apps := []Appointment{}
	events := calendar.Events{}

	addEvents, err := buildDiffLists(apps, &events)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(addEvents))
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

	addEvents, err := buildDiffLists(apps, &events)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(addEvents))
	assert.Equal(t, "sub", addEvents[0].Summary)
	assert.Equal(t, "loc", addEvents[0].Location)
	assert.Equal(t, "body", addEvents[0].Description)
}
