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

	"github.com/abhishek622/moviedock/pkg/discovery"
	"github.com/abhishek622/moviedock/pkg/discovery/consul"
	"github.com/abhishek622/moviedock/user/internal/controller/user"
	httphandler "github.com/abhishek622/moviedock/user/internal/handler/http"
	"github.com/abhishek622/moviedock/user/internal/repository/postgres"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

const serviceName = "user"

func main() {
	var (
		port      = flag.Int("port", 8083, "API handler port")
		consulURL = flag.String("consul-url", "localhost:8500", "Consul URL")
	)
	flag.Parse()
	log.Printf("Starting the movie user service on port %d", port)

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: unable to find .env file")
	}

	// Initialize repository
	repo, err := postgres.New()
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	// Create controller
	ctrl := user.New(repo)

	// Create HTTP handler with Gin
	router := gin.Default()
	handler := httphandler.New(ctrl)
	handler.RegisterRoutes(router)

	// Service discovery setup
	registry, err := consul.NewRegistry(*consulURL)
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	instanceID := discovery.GenerateInstanceID(serviceName)
	serviceAddress := fmt.Sprintf("host.docker.internal:%d", *port)

	// Register service
	if err := registry.Register(ctx, instanceID, serviceName, serviceAddress); err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	// Start health check reporting
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
					log.Printf("Failed to report healthy state: %v", err)
				}
			}
		}
	}()

	// Start HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: router,
	}

	// Graceful shutdown
	var g errgroup.Group
	g.Go(func() error {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("failed to start server: %w", err)
		}
		return nil
	})

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down server...")

	// Deregister service
	if err := registry.Deregister(ctx, instanceID, serviceName); err != nil {
		log.Printf("Failed to deregister service: %v", err)
	}

	// Shutdown server with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	if err := g.Wait(); err != nil {
		log.Printf("Server error: %v", err)
	}

	log.Println("Server stopped")
}
