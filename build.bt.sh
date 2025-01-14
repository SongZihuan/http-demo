#!/bin/sh
# Just use in BT

cd `dirname $0`

time=$(date "+%Y-%m-%d %H:%M:%S")

echo "start run, time: ${time}"

PATH=/usr/local/btgojdk/go1.23.4/bin:$PATH

rm -f ./http-demo

go version && echo "start go mod tidy" && go mod tidy && echo "start go build" && go build -o ./http-demo github.com/SongZihuan/Http-Demo/src/cmd/v1 && echo "build finished"

ls -al ./ | grep http-demo

echo "run finished"