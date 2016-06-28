#!/bin/bash


branch=${1:master}
url=${2:http://localhost:8086}
hardwareID=${3:local}
bench=$(4:.)
timeout=$(5:20m)

cd /root && \
  export GOPATH=$(pwd) && \
  go get -u github.com/sparrc/gdm github.com/influxdata/influxdb github.com/influxdata/grade/cmd/grade && \
  cd $GOPATH/src/github.com/influxdata/influxdb && \
  git checkout $branch && \
  $GOPATH/bin/gdm restore && \
  go test -timeout $timeout -run=^$ -bench=$bench -benchmem ./... > bench.txt && \
  grade \
  -hardwareid="$hardwareID" \
  -influxurl=$url \
  -goversion="$(go version | cut -d' ' -f3-)" \
  -revision="$(git log -1 --format=%H)" \
  -timestamp="$(git log -1 --format=%ct)" \
  < bench.txt
