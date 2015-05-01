echo 'Removing the old version'
ssh pi@calendarpi "rm -rf /home/pi/gopath/src/github.com/sgoertzen/calendarpi/"
echo 'Fetching the new version from github'
ssh pi@calendarpi "export GOPATH=/home/pi/gopath && /home/pi/go/bin/go get github.com/sgoertzen/calendarpi"
echo 'Create the certs directory'
ssh pi@calendarpi "mkdir /home/pi/gopath/src/github.com/sgoertzen/calendarpi/certs"
echo 'Copying the certs into place'
ssh pi@calendarpi "cp /home/pi/gopath/src/github.com/sgoertzen/calendarpi/testsetup/*.pem /home/pi/gopath/src/github.com/sgoertzen/calendarpi/certs"
echo 'Copying our config over'
scp conf.json pi@calendarpi:/home/pi/gopath/src/github.com/sgoertzen/calendarpi
echo 'Building the project'
ssh pi@calendarpi "export GOPATH=/home/pi/gopath && cd /home/pi/gopath/src/github.com/sgoertzen/calendarpi && /home/pi/go/bin/go build"
# TODO: not sure if we want this or just to make it a service.  If a service need to stop it before deploy.
#ssh pi@calendarpi "nohup sudo /home/pi/gopath/src/github.com/sgoertzen/calendarpi"
echo 'Project deployed!'