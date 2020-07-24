#!/bin/sh
enableProxy=$1
proxyIP=$2
tracingPort=$3
metricsPort=$4
distributionPort=$5
flushInterval=$6
applicationName=$7
loadgenFlushInterval=$8
go build -o target/wfgosdktest cmd/beachshirts/main.go
./target/beachshirts conf/shopping.conf "$enableProxy" "$proxyIP" "$tracingPort" "$metricsPort" "$distributionPort" "$flushInterval" "$applicationName"
./target/beachshirts conf/styling.conf "$enableProxy" "$proxyIP" "$tracingPort" "$metricsPort" "$distributionPort" "$flushInterval" "$applicationName"
./target/beachshirts conf/delivery.conf "$enableProxy" "$proxyIP" "$tracingPort" "$metricsPort" "$distributionPort" "$flushInterval" "$applicationName"
./loadgen.sh "$loadgenFlushInterval"


