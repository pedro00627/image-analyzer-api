// Package service contains domain service interfaces
package service

import (
	"context"

	"github.com/pedro00627/image-analyzer-api/internal/domain/entity"
)

// AIAnalyzer defines the interface for AI image analysis services
type AIAnalyzer interface {
	// Analyze analyzes an image and returns tags with confidence scores
	Analyze(ctx context.Context, imageData []byte) (*entity.AnalysisResult, error)
}
