package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/pedro00627/image-analyzer-api/internal/domain/entity"
	"github.com/pedro00627/image-analyzer-api/internal/domain/service/mocks"
)

// TestNewAnalyzeImageUseCase tests the constructor
func TestNewAnalyzeImageUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockValidator := mocks.NewMockImageValidator(ctrl)
	mockAnalyzer := mocks.NewMockAIAnalyzer(ctrl)

	usecase := NewAnalyzeImageUseCase(mockValidator, mockAnalyzer)

	if usecase == nil {
		t.Error("expected non-nil usecase")
	}
}

// TestAnalyzeImageUseCase_ExecuteSuccess tests successful image analysis flow
func TestAnalyzeImageUseCase_ExecuteSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockValidator := mocks.NewMockImageValidator(ctrl)
	mockAnalyzer := mocks.NewMockAIAnalyzer(ctrl)

	imageData := []byte{0xFF, 0xD8, 0xFF} // JPEG header
	filename := "test.jpg"
	mimeType := "image/jpeg"

	// Set up expectations
	mockValidator.EXPECT().
		ValidateType(filename, mimeType).
		Return(nil)

	mockValidator.EXPECT().
		ValidateSize(int64(len(imageData))).
		Return(nil)

	mockValidator.EXPECT().
		ValidateContent(imageData).
		Return(nil)

	expectedResult := entity.NewAnalysisResult([]entity.Tag{
		{Label: "car", Confidence: 0.95},
		{Label: "vehicle", Confidence: 0.87},
	})

	mockAnalyzer.EXPECT().
		Analyze(gomock.Any(), imageData).
		Return(&expectedResult, nil)

	usecase := NewAnalyzeImageUseCase(mockValidator, mockAnalyzer)
	req := AnalyzeImageRequest{
		Filename:  filename,
		ImageData: imageData,
		MimeType:  mimeType,
	}

	result, err := usecase.Execute(context.Background(), req)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if result == nil {
		t.Error("expected non-nil result")
	}

	if len(result.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(result.Tags))
	}
}

// TestAnalyzeImageUseCase_ValidateTypeFails tests when type validation fails
// Note: With new strategy, ValidateType is informational only
// It only fails if content validation fails (bad MIME on bad image)
func TestAnalyzeImageUseCase_ValidateTypeFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockValidator := mocks.NewMockImageValidator(ctrl)
	mockAnalyzer := mocks.NewMockAIAnalyzer(ctrl)

	imageData := []byte{0xFF, 0xD8}
	filename := "test.pdf"
	mimeType := "application/pdf"

	// Size check (should pass)
	mockValidator.EXPECT().
		ValidateSize(int64(len(imageData))).
		Return(nil)

	// Content validation fails (not a valid image)
	mockValidator.EXPECT().
		ValidateContent(imageData).
		Return(errors.New("invalid image content"))

	// ValidateType won't be called because content check fails first

	usecase := NewAnalyzeImageUseCase(mockValidator, mockAnalyzer)
	req := AnalyzeImageRequest{
		Filename:  filename,
		ImageData: imageData,
		MimeType:  mimeType,
	}

	result, err := usecase.Execute(context.Background(), req)

	if err == nil {
		t.Error("expected error from content validation")
	}

	if result != nil {
		t.Error("expected nil result on validation error")
	}
}

// TestAnalyzeImageUseCase_ValidateSizeFails tests when size validation fails
func TestAnalyzeImageUseCase_ValidateSizeFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockValidator := mocks.NewMockImageValidator(ctrl)
	mockAnalyzer := mocks.NewMockAIAnalyzer(ctrl)

	largeImage := make([]byte, 100*1024*1024) // 100MB
	filename := "large.jpg"
	mimeType := "image/jpeg"

	// Size check fails first
	mockValidator.EXPECT().
		ValidateSize(int64(len(largeImage))).
		Return(errors.New("file too large"))

	// Other checks won't be called because size fails first

	usecase := NewAnalyzeImageUseCase(mockValidator, mockAnalyzer)
	req := AnalyzeImageRequest{
		Filename:  filename,
		ImageData: largeImage,
		MimeType:  mimeType,
	}

	result, err := usecase.Execute(context.Background(), req)

	if err == nil {
		t.Error("expected error from size validation")
	}

	if result != nil {
		t.Error("expected nil result on validation error")
	}
}

// TestAnalyzeImageUseCase_ValidateContentFails tests when content validation fails
func TestAnalyzeImageUseCase_ValidateContentFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockValidator := mocks.NewMockImageValidator(ctrl)
	mockAnalyzer := mocks.NewMockAIAnalyzer(ctrl)

	imageData := []byte{0x00, 0x01, 0x02} // Invalid image data
	filename := "test.jpg"
	mimeType := "image/jpeg"

	mockValidator.EXPECT().
		ValidateSize(int64(len(imageData))).
		Return(nil)

	mockValidator.EXPECT().
		ValidateContent(imageData).
		Return(errors.New("invalid image format"))

	// ValidateType won't be called because content check fails first

	usecase := NewAnalyzeImageUseCase(mockValidator, mockAnalyzer)
	req := AnalyzeImageRequest{
		Filename:  filename,
		ImageData: imageData,
		MimeType:  mimeType,
	}

	result, err := usecase.Execute(context.Background(), req)

	if err == nil {
		t.Error("expected error from content validation")
	}

	if result != nil {
		t.Error("expected nil result on validation error")
	}
}

// TestAnalyzeImageUseCase_AnalyzerFails tests when analyzer returns error
func TestAnalyzeImageUseCase_AnalyzerFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockValidator := mocks.NewMockImageValidator(ctrl)
	mockAnalyzer := mocks.NewMockAIAnalyzer(ctrl)

	imageData := []byte{0xFF, 0xD8, 0xFF}
	filename := "test.jpg"
	mimeType := "image/jpeg"

	mockValidator.EXPECT().
		ValidateSize(int64(len(imageData))).
		Return(nil)

	mockValidator.EXPECT().
		ValidateContent(imageData).
		Return(nil)

	mockValidator.EXPECT().
		ValidateType(filename, mimeType).
		Return(nil)

	mockAnalyzer.EXPECT().
		Analyze(gomock.Any(), imageData).
		Return(nil, errors.New("API error"))

	usecase := NewAnalyzeImageUseCase(mockValidator, mockAnalyzer)
	req := AnalyzeImageRequest{
		Filename:  filename,
		ImageData: imageData,
		MimeType:  mimeType,
	}

	result, err := usecase.Execute(context.Background(), req)

	if err == nil {
		t.Error("expected error from analyzer")
	}

	if result != nil {
		t.Error("expected nil result on analyzer error")
	}
}

// TestAnalyzeImageUseCase_ContextCancellation tests behavior with cancelled context
func TestAnalyzeImageUseCase_ContextCancellation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockValidator := mocks.NewMockImageValidator(ctrl)
	mockAnalyzer := mocks.NewMockAIAnalyzer(ctrl)

	imageData := []byte{0xFF, 0xD8, 0xFF}
	filename := "test.jpg"
	mimeType := "image/jpeg"

	// Set up validations
	mockValidator.EXPECT().
		ValidateType(filename, mimeType).
		Return(nil)

	mockValidator.EXPECT().
		ValidateSize(int64(len(imageData))).
		Return(nil)

	mockValidator.EXPECT().
		ValidateContent(imageData).
		Return(nil)

	// Analyzer should still receive the call with cancelled context
	mockAnalyzer.EXPECT().
		Analyze(gomock.Any(), imageData).
		Return(nil, context.Canceled)

	usecase := NewAnalyzeImageUseCase(mockValidator, mockAnalyzer)
	req := AnalyzeImageRequest{
		Filename:  filename,
		ImageData: imageData,
		MimeType:  mimeType,
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err := usecase.Execute(ctx, req)

	if err == nil {
		t.Error("expected error from cancelled context")
	}

	if result != nil {
		t.Error("expected nil result on context error")
	}
}

// TestAnalyzeImageUseCase_EmptyImageData tests with empty image data
func TestAnalyzeImageUseCase_EmptyImageData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockValidator := mocks.NewMockImageValidator(ctrl)
	mockAnalyzer := mocks.NewMockAIAnalyzer(ctrl)

	imageData := []byte{}
	filename := "empty.jpg"
	mimeType := "image/jpeg"

	mockValidator.EXPECT().
		ValidateSize(int64(len(imageData))).
		Return(errors.New("image is empty"))

	usecase := NewAnalyzeImageUseCase(mockValidator, mockAnalyzer)
	req := AnalyzeImageRequest{
		Filename:  filename,
		ImageData: imageData,
		MimeType:  mimeType,
	}

	result, err := usecase.Execute(context.Background(), req)

	if err == nil {
		t.Error("expected error for empty image")
	}

	if result != nil {
		t.Error("expected nil result on validation error")
	}
}

// TestAnalyzeImageUseCase_MultipleValidationSteps tests the order of validations
// Order: size → content → type (size first for quick rejection, type last as informational)
func TestAnalyzeImageUseCase_MultipleValidationSteps(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockValidator := mocks.NewMockImageValidator(ctrl)
	mockAnalyzer := mocks.NewMockAIAnalyzer(ctrl)

	imageData := []byte{0xFF, 0xD8, 0xFF}
	filename := "test.png"
	mimeType := "image/png"

	gomock.InOrder(
		mockValidator.EXPECT().ValidateSize(int64(len(imageData))).Return(nil),
		mockValidator.EXPECT().ValidateContent(imageData).Return(nil),
		mockValidator.EXPECT().ValidateType(filename, mimeType).Return(nil),
		mockAnalyzer.EXPECT().Analyze(gomock.Any(), imageData).Return(
			&entity.AnalysisResult{Tags: []entity.Tag{{Label: "test", Confidence: 0.9}}},
			nil,
		),
	)

	usecase := NewAnalyzeImageUseCase(mockValidator, mockAnalyzer)
	req := AnalyzeImageRequest{
		Filename:  filename,
		ImageData: imageData,
		MimeType:  mimeType,
	}

	result, err := usecase.Execute(context.Background(), req)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if result == nil {
		t.Error("expected non-nil result")
	}
}
