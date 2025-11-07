package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/abhishek622/moviedock/pkg/discovery"
	"github.com/abhishek622/moviedock/pkg/discovery/consul"
)

const serviceName = "user"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()
	log.Printf("Starting the movie user service on port %d", port)
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Printf("Failed to report healthy state: %v", err.Error())
			}
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)
}
