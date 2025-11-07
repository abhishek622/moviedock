package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/abhishek622/moviedock/pkg/discovery"
	"github.com/abhishek622/moviedock/pkg/discovery/consul"
)

const serviceName = "metadata"

func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "API handler port")
	flag.Parse()
	log.Printf("Starting the movie metadata service on port %d", port)
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
	// repo := postgres.New()
	// svc := metadata.New(repo)
	// h := httphandler.New(svc)
	// http.Handle("/metadata", http.HandleFunc(h.GetMetadataByID))
	// if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
	// 	panic(err)
	// }

}
