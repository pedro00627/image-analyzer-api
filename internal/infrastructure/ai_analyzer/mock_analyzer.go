package ai_analyzer

import (
	"context"

	"github.com/pedro00627/image-analyzer-api/internal/domain/entity"
)

// MockAnalyzer is a mock implementation for testing without Google Vision API
type MockAnalyzer struct {
	// For testing: can be configured to return specific results or errors
}

var _ Analyzer = (*MockAnalyzer)(nil)

// NewMockAnalyzer creates a new mock analyzer for testing
func NewMockAnalyzer() *MockAnalyzer {
	return &MockAnalyzer{}
}

// Analyze returns mock analysis results
func (m *MockAnalyzer) Analyze(ctx context.Context, imageData []byte) (*entity.AnalysisResult, error) {
	// Return mock results for testing
	result := entity.NewAnalysisResult([]entity.Tag{
		{Label: "mock-tag-1", Confidence: 0.95},
		{Label: "mock-tag-2", Confidence: 0.87},
		{Label: "test", Confidence: 0.75},
	})
	return &result, nil
}

// Close is a no-op for mock
func (m *MockAnalyzer) Close() error {
	return nil
}
