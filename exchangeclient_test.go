package main

import (
	"github.com/stretchr/testify/assert"
	"strconv"
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
	// Looks somethign like <mes:CalendarView MaxEntriesReturned="100" StartDate="2015-04-21T05:59:57Z" EndDate="2015-05-05T05:59:57Z"/>
	calendarline := request[start:end]
	keyvaluepairs := strings.Split(calendarline, " ")

	// Verify the dates are there and the max entries contains a number
	for _, keyvalue := range keyvaluepairs {
		if strings.Index(keyvalue, "=") > 0 {
			parts := strings.Split(keyvalue, "=")
			assert.NotNil(t, parts)
			switch parts[0] {
			case "MaxEntriesReturned":
				numstring := parts[1][1 : len(parts[1])-1]
				i, err := strconv.ParseInt(numstring, 0, 64)
				assert.Nil(t, err)
				assert.True(t, i > 0)
			case "StartDate":
				assert.Equal(t, 22, len(parts[1]))
			case "EndDate":
				assert.Equal(t, 25, len(parts[1])) // Length includes ending xml />
			}
		}
	}
}
