#!/bin/sh

# Make sure we're in the right folder
cd /home/user/app

# Remove all previous executables
rm ./bin/*.gin

# Set path to config file
export CONFIG_FILE=/config.yaml

# Build will be slow for the first time
# Put a reassuring message to the developers
(sleep 8 && echo 'Building ohana, please wait... (might take a while)' > /dev/stdout) &

# Grab dependencies
go mod download -x

# Start gin
exec /go/bin/gin \
  --immediate \
  --port $GIN_PORT \
  --appPort $APP_PORT \
  --build cmd/ohana/ \
  --bin ./bin/$(hostname).gin \
  --buildArgs "-tags osusergo,netgo"
