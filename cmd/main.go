package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pedro00627/image-analyzer-api/internal/api/server"
	"github.com/pedro00627/image-analyzer-api/internal/infrastructure/bootstrap"
)

func main() {
	// Set Gin mode
	gin.SetMode(gin.DebugMode)

	ctx := context.Background()

	// Bootstrap dependencies
	container, err := bootstrap.Bootstrap(ctx)
	if err != nil {
		log.Printf("Failed to bootstrap dependencies: %v", err)
		os.Exit(1)
	}

	// Create and start HTTP server
	srv := server.New(container)
	if err := srv.Start(); err != nil {
		log.Printf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
