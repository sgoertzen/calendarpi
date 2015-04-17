# calendarpi
CalendarPi will sync calendar appointments between an Exchange Server and Google Calendar.  It is made to run on the RaspberryPi but can be run on any machine that supports Go.

# Setup Instructions
Download and install Go from https://golang.org/dl/ or from source by running:
git clone https://go.googlesource.com/go
cd go
git checkout go1.4
cd src
./all.bash

Download and install Mecurial by running:
sudo apt-get install mercurial

Run the following commands in terminal
go get google.golang.org/api/calendar/v3
go get golang.org/x/oauth2

To build:
go build

To run:
./calendarpi
