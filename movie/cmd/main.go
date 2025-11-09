package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abhishek622/moviedock/movie/internal/controller/movie"
	metadatagateway "github.com/abhishek622/moviedock/movie/internal/gateway/metadata/http"
	ratinggateway "github.com/abhishek622/moviedock/movie/internal/gateway/rating/http"
	httphandler "github.com/abhishek622/moviedock/movie/internal/handler/http"
	"github.com/abhishek622/moviedock/pkg/discovery"
	"github.com/abhishek622/moviedock/pkg/discovery/consul"
	"github.com/gin-gonic/gin"
)

const serviceName = "movie"

func main() {
	var port int
	flag.IntVar(&port, "port", 8084, "API handler port")
	flag.Parse()

	// Initialize service discovery
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		log.Fatalf("Failed to create service registry: %v", err)
	}

	// Initialize gateways
	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)

	// Initialize service and controller
	svc := movie.New(ratingGateway, metadataGateway)

	// Create Gin router
	router := gin.Default()

	// Initialize and register movie handler
	handler := httphandler.New(svc)
	handler.RegisterRoutes(router)

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	// Register service with service discovery
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("Starting movie service on port %d", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Start health reporting
	healthTicker := time.NewTicker(1 * time.Second)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Health reporting goroutine
	go func() {
		for range healthTicker.C {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Printf("Failed to report healthy state: %v", err)
			} else {
				log.Println("Successfully reported healthy state")
			}
		}
	}()

	// Wait for interrupt signal
	<-quit
	healthTicker.Stop()
	log.Println("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Deregister from service discovery
	if err := registry.Deregister(ctx, instanceID, serviceName); err != nil {
		log.Printf("Error deregistering service: %v", err)
	}

	log.Println("Server exiting")
}
