package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestBuildCalendarRequest(t *testing.T) {
	requestbytes := buildCalendarRequest("black", "ninja")
	request := string(requestbytes)
	assert.NotNil(t, request)

	// Only testing the two lines that get edited
	assert.True(t, strings.Contains(request, `<typ:FolderId Id="black" ChangeKey="ninja" />`))

	// Date string should always be the same length so this should always be the same
	start := strings.Index(request, "<mes:CalendarView")
	end := strings.Index(request, "<mes:ParentFolderIds")
	calendarline := request[start:end]
	// Looks somethign like <mes:CalendarView MaxEntriesReturned="5" StartDate="2015-04-21T05:59:57Z" EndDate="2015-05-05T05:59:57Z"/>
	assert.Equal(t, 116, len(calendarline))
}
