package service

import (
	"context"

	"github.com/pedro00627/image-analyzer-api/internal/application/usecase"
	"github.com/pedro00627/image-analyzer-api/internal/domain/entity"
	"github.com/pedro00627/image-analyzer-api/internal/domain/error"
)

// ImageAnalysisService defines the contract for image analysis operations
type ImageAnalysisService interface {
	AnalyzeImage(ctx context.Context, request *AnalyzeImageRequest) (*entity.AnalysisResult, *error.DomainError)
}

// AnalyzeImageRequest represents a request to analyze an image
type AnalyzeImageRequest struct {
	Filename  string
	MimeType  string
	ImageData []byte
}

// imageAnalysisServiceImpl implements ImageAnalysisService
type imageAnalysisServiceImpl struct {
	useCase usecase.AnalyzeImageUseCaseInterface
}

// NewImageAnalysisService creates a new ImageAnalysisService
func NewImageAnalysisService(uc usecase.AnalyzeImageUseCaseInterface) ImageAnalysisService {
	return &imageAnalysisServiceImpl{
		useCase: uc,
	}
}

// AnalyzeImage analyzes an image and returns the analysis result
func (s *imageAnalysisServiceImpl) AnalyzeImage(ctx context.Context, request *AnalyzeImageRequest) (*entity.AnalysisResult, *error.DomainError) {
	// Delegate to usecase
	useCaseReq := usecase.AnalyzeImageRequest{
		Filename:  request.Filename,
		MimeType:  request.MimeType,
		ImageData: request.ImageData,
	}

	return s.useCase.Execute(ctx, useCaseReq)
}
