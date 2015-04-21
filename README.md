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

## First run:
The first time the software runs it will error out and create a conf.json file in the root directory.  You will need to edit this file and put in proper values.

## Usage
Once the software is running point your browser to https://yourmachinename:10443/ (Or whatever port you put into the config file)
