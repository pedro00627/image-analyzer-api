package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pedro00627/image-analyzer-api/internal/api/handler"
	"github.com/pedro00627/image-analyzer-api/internal/application/service"
	"github.com/pedro00627/image-analyzer-api/internal/infrastructure/bootstrap"
)

// SetupRoutes configures all API routes
func SetupRoutes(engine *gin.Engine, container bootstrap.DependencyContainer) {
	// Configure CORS
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://pedro00627.com", "https://www.pedro00627.com", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Get use case from container
	useCase := container.GetAnalyzeImageUseCase()

	// Create service layer
	analysisService := service.NewImageAnalysisService(useCase)

	// Create handlers
	imageHandler := handler.NewImageHandler(analysisService)

	// Register routes
	api := engine.Group("/api")
	{
		api.POST("/analyze", imageHandler.AnalyzeImage)
	}
}
