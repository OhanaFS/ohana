#!/bin/sh

# Make sure we're in the right folder
cd /home/user/app

# Create the binary
export bin_name=$(mktemp -p ./bin -u)

# Set path to config file
export CONFIG_FILE=/config.yaml

# Build will be slow for the first time
# Put a reassuring message to the developers
(sleep 8 && echo 'Building ohana, please wait... (usually takes one minute)' > /dev/stdout) &
(sleep 9 && echo "Building to $bin_name" > /dev/stdout) &

# Grab dependencies
go mod download -x

# Start gin
exec /go/bin/gin \
  --immediate \
  --port "$GIN_PORT" \
  --appPort "$APP_PORT" \
  --build cmd/ohana/ \
  --bin "$bin_name" \
  --buildArgs "-tags osusergo,netgo"
