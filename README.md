# CalendarPi
CalendarPi will sync calendar appointments between an Exchange Server and Google Calendar.  It is made to run on the RaspberryPi but can be run on any machine that supports Go.

[![Build Status](https://travis-ci.org/sgoertzen/calendarpi.svg?branch=master)](https://travis-ci.org/sgoertzen/calendarpi)
[![Codacy Badge](https://www.codacy.com/project/badge/f0dedbbcb471499eb47456cf954018d3)](https://www.codacy.com/app/sgoertzen/calendarpi)
[![Coverage Status](https://coveralls.io/repos/sgoertzen/calendarpi/badge.svg)](https://coveralls.io/r/sgoertzen/calendarpi)
[![Author Badge](https://img.shields.io/badge/awesome-totally-green.svg)](https://github.com/sgoertzen)

CalendarPi is a small website that allows users to sync their exchange calenders with google calendar.  Multiple users can all use the software and it will maintain each calenadr separately.  

![Screenshot](https://github.com/sgoertzen/calendarpi/blob/master/images/ScreenShot.png)

### Details
CalendarPi supports one directional sync from Exchange to Google Calendar.  It will automatically add new items, update exisiting items and delete removed items.

### Security
CalendarPi secures all credentials by storing them in an encrypted file locally.  The master password for encrypting and decrypting is never stored on the system.  It must be entered each time the software starts up and is just kept in memory.  This prohibts anyone who has physical access to the machine from recovering any sensitive data.

## Setup Instructions
Download and install Go from https://golang.org/dl/
Set your GOPATH variable to whatever you want
```sh
export GOPATH=~/GoPath
```
Install calendarpi by running
```sh
go get github.com/sgoertzen/calendarpi
```
Install all dependencies by running
```sh
go get -t ./...
```
Finally switch to the code directory
```sh
cd $GOPATH/src/github.com/sgoertzen/calendarpi
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

## Installing Go from source
If you are so inclined you can install go from source
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
go get github.com/kennygrant/sanitize
```
