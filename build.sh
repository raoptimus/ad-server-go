#!/bin/sh
export GOPATH=$(pwd);
echo $GOPATH;
rm -rf bin/*
go build -o bin/TeaserNet "./src/tc/Advertiser/TeaserNet" && chmod 0755 bin/TeaserNet && echo "\nTeaserNet builded" >&2 &
go build -o bin/TubeContext "./src/tc/Advertiser/TubeContext" && chmod 0755 bin/TubeContext && echo "\nTubeContext builded" >&2 &
go build -o bin/EroAdvertising "./src/tc/Advertiser/EroAdvertising" && chmod 0755 bin/EroAdvertising && echo "\nEroAdvertising builded" >&2 &
go build -o bin/Server "./src/tc/Publisher/Server" && chmod 0755 bin/Server && echo "\nServer builded" >&2 &
go build -o bin/Click "./src/tc/Publisher/Click" && chmod 0755 bin/Click && echo "\nClick builded" >&2 &
go build -o bin/RTBWin "./src/tc/Advertiser/RTBWin" && chmod 0755 bin/RTBWin && echo "\nRTBWin builded" >&2 &



