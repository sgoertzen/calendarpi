# CalendarPi
CalendarPi will sync calendar appointments between an Exchange Server and Google Calendar.  It is made to run on the RaspberryPi but can be run on any machine that supports Go.

[![Build Status](https://travis-ci.org/sgoertzen/calendarpi.svg?branch=master)](https://travis-ci.org/sgoertzen/calendarpi)
[![Codacy Badge](https://www.codacy.com/project/badge/f0dedbbcb471499eb47456cf954018d3)](https://www.codacy.com/app/sgoertzen/calendarpi)

### Version
0.1.0

This software is not fully functional but is actively being developed.  

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
go get github.com/stretchr/testify/assert
```

## To build:
```sh
go build
```

## To test:
```sh
go test
```

## To run:
```sh
./calendarpi
```

## First run:
The first time the software runs it will error out and create a conf.json file in the root directory.  You will need to edit this file and put in proper values.

## Usage
Once the software is running point your browser to https://yourmachinename:10443/ (Or whatever port you put into the config file)
