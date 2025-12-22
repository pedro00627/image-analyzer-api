package service

import (
	"context"

	"github.com/pedro00627/image-analyzer-api/internal/domain/entity"
)

//go:generate mockgen -source=./ai_analyzer.go -destination=mocks/mock_ai_analyzer.go -package=mocks

// AIAnalyzer defines the interface for AI image analysis services
type AIAnalyzer interface {
	// Analyze analyzes an image and returns tags with confidence scores
	Analyze(ctx context.Context, imageData []byte) (*entity.AnalysisResult, error)
}
