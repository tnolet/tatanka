#!/bin/bash
set -x
# build the app for linux/i386 an build the docker container
GOOS=linux GOARCH=386 go build
mv tatanka target/linux_i386/
docker build -t tnolet/tatanka:$1 .