# This will deploy the code to a raspberrypi device called "calendarpi"
echo 'Stopping service'
ssh pi@calendarpi "sudo killall calendarpi"
echo 'Backing up users data'
ssh pi@calendarpi "mv /home/pi/gopath/src/github.com/sgoertzen/calendarpi/blabbersnatzle.bak /home/pi/gopath/src/github.com/sgoertzen/"
echo 'Removing the old version'
ssh pi@calendarpi "rm -rf /home/pi/gopath/src/github.com/sgoertzen/calendarpi/"
ssh pi@calendarpi "rm -rf /home/pi/gopath/src/github.com/sgoertzen/html2text/"
ssh pi@calendarpi "rm -rf /home/pi/gopath/src/github.com/sgoertzen/xchango/"
echo 'Fetching the new version from github'
ssh pi@calendarpi "export GOPATH=/home/pi/gopath && /home/pi/go/bin/go get github.com/sgoertzen/calendarpi"
echo 'Create the certs directory'
ssh pi@calendarpi "mkdir /home/pi/gopath/src/github.com/sgoertzen/calendarpi/certs"
echo 'Copying the certs into place'
scp certs/cert.pem certs/key.pem pi@calendarpi:/home/pi/gopath/src/github.com/sgoertzen/calendarpi/certs
echo 'Copying our config over'
scp conf.json pi@calendarpi:/home/pi/gopath/src/github.com/sgoertzen/calendarpi
echo 'Fetching dependencies'
ssh pi@calendarpi "export GOPATH=/home/pi/gopath && cd /home/pi/gopath/src/github.com/sgoertzen/calendarpi && /home/pi/go/bin/go get ./..."
echo 'Building the project'
ssh pi@calendarpi "export GOPATH=/home/pi/gopath && cd /home/pi/gopath/src/github.com/sgoertzen/calendarpi && /home/pi/go/bin/go build"
echo 'Restoring users data'
ssh pi@calendarpi "mv /home/pi/gopath/src/github.com/sgoertzen/blabbersnatzle.bak /home/pi/gopath/src/github.com/sgoertzen/calendarpi"
echo 'Restarting service'
ssh pi@calendarpi "cd /home/pi/gopath/src/github.com/sgoertzen/calendarpi && sudo /home/pi/gopath/src/github.com/sgoertzen/calendarpi/calendarpi &" &
echo 'Project deployed!'