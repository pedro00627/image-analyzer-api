package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pedro00627/image-analyzer-api/internal/api/handler"
	"github.com/pedro00627/image-analyzer-api/internal/api/middleware"
	"github.com/pedro00627/image-analyzer-api/internal/application/service"
	"github.com/pedro00627/image-analyzer-api/internal/infrastructure/bootstrap"
	"golang.org/x/time/rate"
)

// SetupRoutes configures all API routes
func SetupRoutes(engine *gin.Engine, container bootstrap.DependencyContainer) {
	config := container.GetConfig()

	// Configure CORS from environment
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     config.GetAllowedOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Configure Rate Limiting from config
	perMinute := config.GetRateLimitPerMinute()
	burst := config.GetRateLimitBurst()
	rateLimiter := middleware.NewRateLimiter(rate.Every(time.Minute/time.Duration(perMinute)), burst)
	engine.Use(rateLimiter.Limit())

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
