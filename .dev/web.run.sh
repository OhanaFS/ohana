#!/bin/sh

# Enter the web directory
cd /home/user/app/web

# Download dependencies
yarn install

# Run develpoment server
exec yarn dev --host
