#/bin/bash

set -e

echo "running tests"
go test ./...

cd bot/

echo "compiling for 386"
GOOS=linux GOARCH=386 go build

echo "removing old distribution"
rm -f foosbot.zip

echo "packaging new distribution"
zip foosbot.zip bot

echo "deploying new distribution"
rsync --progress foosbot.zip ubuntu@mubit.io:~/foosbot

echo "extracting new distribution on remote"
ssh ubuntu@mubit.io "cd foosbot/ && unzip foosbot.zip"

