#!/bin/sh
export GOPATH=$(pwd);
echo $GOPATH;
rm -rf bin/Server
go build -o bin/Server "./src/tc/Publisher/Server" && chmod 0755 bin/Server
bin/Server

