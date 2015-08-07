#!/bin/bash
set -x
# build the app for linux/i386
GOOS=linux GOARCH=386 go build 
mv tatanka target/linux_i386/
