# CalendarPi
CalendarPi will sync calendar appointments between an Exchange Server and Google Calendar.  It is made to run on the RaspberryPi but can be run on any machine that supports Go.

This software is not fully functional but is actively being developed.  

### Version
0.1.0

## Setup Instructions
Download and install Go from https://golang.org/dl/ or from source by running:
```sh
git clone https://go.googlesource.com/go
cd go
git checkout go1.4
cd src
./all.bash
```

Download and install Mecurial by running:
```sh
sudo apt-get install mercurial
```

Run the following commands in terminal
```sh
go get google.golang.org/api/calendar/v3
go get golang.org/x/oauth2
```

## To build:
```sh
go build
```

## To run:
```sh
./calendarpi
```


