package bootstrap

import (
	"context"

	"github.com/pedro00627/image-analyzer-api/internal/application/usecase"
	"github.com/pedro00627/image-analyzer-api/internal/domain/service"
	"github.com/pedro00627/image-analyzer-api/internal/infrastructure/ai_analyzer"
	"github.com/pedro00627/image-analyzer-api/internal/infrastructure/config"
	"github.com/pedro00627/image-analyzer-api/internal/infrastructure/validator"
)

// DependencyContainer defines the contract for accessing dependencies
type DependencyContainer interface {
	GetConfig() config.Config
	GetAnalyzeImageUseCase() usecase.AnalyzeImageUseCaseInterface
}

// Container holds all application dependencies
type Container struct {
	cfg     config.Config
	useCase usecase.AnalyzeImageUseCaseInterface
}

// Bootstrap initializes all dependencies and returns a DependencyContainer
func Bootstrap(ctx context.Context) (DependencyContainer, error) {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	return BootstrapWithConfig(ctx, cfg)
}

// BootstrapWithConfig initializes all dependencies with a given config (useful for testing)
func BootstrapWithConfig(ctx context.Context, cfg config.Config) (DependencyContainer, error) {
	// Initialize AI analyzer
	aiAnalyzer := initializeAIAnalyzer(ctx, cfg)

	// Initialize validator
	imgValidator := validator.NewImageValidator(cfg)

	// Initialize use case
	uc := usecase.NewAnalyzeImageUseCase(imgValidator, aiAnalyzer)

	return &Container{
		cfg:     cfg,
		useCase: uc,
	}, nil
}

// GetConfig returns the application configuration
func (c *Container) GetConfig() config.Config {
	return c.cfg
}

// GetAnalyzeImageUseCase returns the AnalyzeImageUseCase
func (c *Container) GetAnalyzeImageUseCase() usecase.AnalyzeImageUseCaseInterface {
	return c.useCase
}

// initializeAIAnalyzer initializes the AI analyzer with fallback to mock
func initializeAIAnalyzer(ctx context.Context, cfg config.Config) service.AIAnalyzer {
	analyzer, err := ai_analyzer.NewGoogleVisionAnalyzer(ctx, cfg)
	if err != nil {
		// Fallback to mock analyzer if Google Vision fails
		return ai_analyzer.NewMockAnalyzer()
	}
	return analyzer
}
