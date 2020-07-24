#!/bin/sh
enableProxy=$1
cluster=$2
token=$3
batchSize=$4
maxBufferSize=$5
flushInterval=$6
applicationName=$7
loadgenFlushInterval=$8
go build -o target/wfgosdktest cmd/beachshirts/main.go
./target/beachshirts conf/shopping.conf "$enableProxy" "$cluster" "$token" "$batchSize" "$maxBufferSize" "$flushInterval" "$applicationName"
./target/beachshirts conf/styling.conf "$enableProxy" "$cluster" "$token" "$batchSize" "$maxBufferSize" "$flushInterval" "$applicationName"
./target/beachshirts conf/delivery.conf "$enableProxy" "$cluster" "$token" "$batchSize" "$maxBufferSize" "$flushInterval" "$applicationName"
./loadgen.sh "$loadgenFlushInterval"


