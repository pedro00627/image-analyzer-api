package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pedro00627/image-analyzer-api/internal/api/router"
	"github.com/pedro00627/image-analyzer-api/internal/infrastructure/bootstrap"
)

// HTTPServer defines the contract for HTTP server operations
type HTTPServer interface {
	Start() error
}

// Server implements the HTTPServer interface
type Server struct {
	engine *gin.Engine
	config bootstrap.DependencyContainer
	port   int
}

// New creates a new HTTPServer instance
func New(container bootstrap.DependencyContainer) HTTPServer {
	engine := gin.Default()

	// Health check endpoint
	engine.GET("/health", healthCheck)

	// Setup API routes
	router.SetupRoutes(engine, container)

	port := 8080

	return &Server{
		engine: engine,
		config: container,
		port:   port,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Starting server on %s", addr)
	return s.engine.Run(addr)
}

// healthCheck handles the health check endpoint
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}
