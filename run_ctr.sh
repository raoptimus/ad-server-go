#!/bin/sh
export GOPATH=$(pwd);
echo $GOPATH;
rm -rf bin/Ctr
go build -o bin/Ctr "./src/tc/Ctr" && chmod 0755 bin/Ctr
bin/Ctr

