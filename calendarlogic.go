package main

import (
	"fmt"
	"google.golang.org/api/calendar/v3"
	"log"
	"time"
)

func processAppointments(user User, apps []Appointment) {
	for _, app := range apps {
		log.Println(app.Subject)

		client := getClient(user)
		srv, err := calendar.New(client)
		if err != nil {
			// TODO: switch failures to Fatals
			log.Fatalf("Unable to retrieve calendar Client %v", err)
		}

		// TODO: This is just a test for now.
		// Need to actually implement the logic
		t := time.Now().Format(time.RFC3339)
		events, err := srv.Events.List("primary").ShowDeleted(false).
			SingleEvents(true).
			TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve next ten of the user's events. %v", err)
		}

		fmt.Println("Upcoming events:")
		if len(events.Items) > 0 {
			for _, i := range events.Items {
				var when string
				// If the DateTime is an empty string the Event is an all-day Event.
				// So only Date is available.
				if i.Start.DateTime != "" {
					when = i.Start.DateTime
				} else {
					when = i.Start.Date
				}
				fmt.Printf("%s (%s)\n", i.Summary, when)
			}
		} else {
			fmt.Printf("No upcoming events found.")
		}
	}
}
