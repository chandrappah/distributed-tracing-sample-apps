package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"os"
	"strconv"
	"io"

	. "wavefront.com/hackathon/beachshirts/internal"
	"wavefront.com/hackathon/beachshirts/services/delivery"
	"wavefront.com/hackathon/beachshirts/services/shopping"
	"wavefront.com/hackathon/beachshirts/services/styling"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: beachshirts <config_file>\n")
		os.Exit(1)
	}

	InitGlobalConfig()

	file, ferr := os.Open(os.Args[1])
	if ferr != nil {
		fmt.Println(ferr)
		os.Exit(2)
	}
	if derr := json.NewDecoder(file).Decode(&GlobalConfig); derr != nil {
		fmt.Println(derr)
		os.Exit(2)
	}

	var server Server
	var closer io.Closer
	serviceName := GlobalConfig.Service 

	if strings.Contains(os.Args[2], "true") {
		proxyIP := os.Args[3]
		tracingPort, _ := strconv.Atoi(os.Args[4])
		metricsPort, _ := strconv.Atoi(os.Args[5])
		distributedPort, _ := strconv.Atoi(os.Args[6])
		flushInterval, _ := strconv.Atoi(os.Args[7])
		applicationName := os.Args[8]
		closer = NewProxyGlobalTracer(serviceName + "-test", proxyIP, tracingPort, metricsPort, distributedPort, flushInterval, applicationName)
	} else {
		cluster := os.Args[3]
		token := os.Args[4]
		batchSize, _ := strconv.Atoi(os.Args[5])
                maxBufferSize, _ := strconv.Atoi(os.Args[6])
                flushInterval, _ := strconv.Atoi(os.Args[7])
                applicationName := os.Args[8]
		closer = NewDirectGlobalTracer(serviceName + "-test", cluster, token, batchSize, maxBufferSize, flushInterval, applicationName)

}

	defer closer.Close()

	switch GlobalConfig.Service {
	case "shopping":
		server = shopping.NewServer()

	case "styling":
		server = styling.NewServer()

	case "delivery":
		server = delivery.NewServer()

	default:
		fmt.Printf("Unrecognized beachshirts service: %s\n", os.Args[1])
		os.Exit(1)
	}

	if serr := server.Start(); serr != nil {
		fmt.Println(serr.Error())
		os.Exit(1)
	}
}

