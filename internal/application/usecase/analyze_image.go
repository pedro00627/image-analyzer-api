package usecase

import (
	"context"
	"fmt"

	"github.com/pedro00627/image-analyzer-api/internal/domain/entity"
	"github.com/pedro00627/image-analyzer-api/internal/domain/service"
)

//go:generate mockgen -source=./analyze_image.go -destination=mocks/mock_analyze_image_usecase.go -package=mocks

// AnalyzeImageUseCaseInterface defines the contract for image analysis orchestration
type AnalyzeImageUseCaseInterface interface {
	Execute(ctx context.Context, req AnalyzeImageRequest) (*entity.AnalysisResult, error)
}

// AnalyzeImageUseCase orchestrates image validation and analysis
type AnalyzeImageUseCase struct {
	validator service.ImageValidator
	analyzer  service.AIAnalyzer
}

var _ AnalyzeImageUseCaseInterface = (*AnalyzeImageUseCase)(nil)

// NewAnalyzeImageUseCase creates a new AnalyzeImageUseCase with dependency injection
func NewAnalyzeImageUseCase(validator service.ImageValidator, analyzer service.AIAnalyzer) *AnalyzeImageUseCase {
	return &AnalyzeImageUseCase{
		validator: validator,
		analyzer:  analyzer,
	}
}

// AnalyzeImageRequest represents the input for image analysis
type AnalyzeImageRequest struct {
	Filename  string
	MimeType  string
	ImageData []byte
}

// Execute orchestrates the image validation and analysis workflow
// Returns an AnalysisResult or error if validation fails
func (u *AnalyzeImageUseCase) Execute(ctx context.Context, req AnalyzeImageRequest) (*entity.AnalysisResult, error) {
	// Validate image type
	if err := u.validator.ValidateType(req.Filename, req.MimeType); err != nil {
		return nil, fmt.Errorf("type validation failed: %w", err)
	}

	// Validate image size
	if err := u.validator.ValidateSize(int64(len(req.ImageData))); err != nil {
		return nil, fmt.Errorf("size validation failed: %w", err)
	}

	// Validate image content
	if err := u.validator.ValidateContent(req.ImageData); err != nil {
		return nil, fmt.Errorf("content validation failed: %w", err)
	}

	// Analyze image
	result, err := u.analyzer.Analyze(ctx, req.ImageData)
	if err != nil {
		return nil, fmt.Errorf("analysis failed: %w", err)
	}

	return result, nil
}
