package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pedro00627/image-analyzer-api/internal/api/handler"
	"github.com/pedro00627/image-analyzer-api/internal/application/service"
	"github.com/pedro00627/image-analyzer-api/internal/infrastructure/bootstrap"
)

// SetupRoutes configures all API routes
func SetupRoutes(engine *gin.Engine, container bootstrap.DependencyContainer) {
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
