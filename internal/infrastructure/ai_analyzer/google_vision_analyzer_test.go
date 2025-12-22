package ai_analyzer

import (
	"context"
	"errors"
	"testing"

	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"github.com/golang/mock/gomock"

	"github.com/pedro00627/image-analyzer-api/internal/infrastructure/ai_analyzer/mocks"
	cfgmocks "github.com/pedro00627/image-analyzer-api/internal/infrastructure/config/mocks"
)

// TestGoogleVisionAnalyzer_AnalyzeEmptyImage tests error handling for empty/nil image data
func TestGoogleVisionAnalyzer_AnalyzeEmptyImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCfg := cfgmocks.NewMockConfig(ctrl)
	mockClient := mocks.NewMockVisionClient(ctrl)

	analyzer := NewGoogleVisionAnalyzerWithClient(mockClient, mockCfg)

	// Test nil image
	_, err := analyzer.Analyze(context.Background(), nil)
	if err == nil {
		t.Error("expected error for nil image data, got nil")
	}

	// Test empty image
	_, err = analyzer.Analyze(context.Background(), []byte{})
	if err == nil {
		t.Error("expected error for empty image data, got nil")
	}
}

// TestGoogleVisionAnalyzer_NewGoogleVisionAnalyzer tests initialization without credentials
func TestGoogleVisionAnalyzer_NewGoogleVisionAnalyzer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCfg := cfgmocks.NewMockConfig(ctrl)

	// This will fail without proper Google Cloud credentials
	_, err := NewGoogleVisionAnalyzer(context.Background(), mockCfg)
	if err == nil {
		t.Error("expected error without proper Google credentials")
	}
}

// TestGoogleVisionAnalyzer_AnalyzeSuccess tests successful image analysis
func TestGoogleVisionAnalyzer_AnalyzeSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCfg := cfgmocks.NewMockConfig(ctrl)
	mockClient := mocks.NewMockVisionClient(ctrl)

	imageData := []byte("test image data")

	// Set up mock response
	expectedResponse := &visionpb.BatchAnnotateImagesResponse{
		Responses: []*visionpb.AnnotateImageResponse{
			{
				LabelAnnotations: []*visionpb.EntityAnnotation{
					{
						Description: "label1",
						Score:       0.95,
					},
					{
						Description: "label2",
						Score:       0.87,
					},
				},
			},
		},
	}

	// Expect BatchAnnotateImages to be called
	mockClient.EXPECT().
		BatchAnnotateImages(gomock.Any(), gomock.Any()).
		Return(expectedResponse, nil)

	analyzer := NewGoogleVisionAnalyzerWithClient(mockClient, mockCfg)

	result, err := analyzer.Analyze(context.Background(), imageData)
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

// TestGoogleVisionAnalyzer_AnalyzeAPIError tests API call error handling
func TestGoogleVisionAnalyzer_AnalyzeAPIError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCfg := cfgmocks.NewMockConfig(ctrl)
	mockClient := mocks.NewMockVisionClient(ctrl)

	imageData := []byte("test image data")

	// Set up mock to return error
	mockClient.EXPECT().
		BatchAnnotateImages(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("API error"))

	analyzer := NewGoogleVisionAnalyzerWithClient(mockClient, mockCfg)

	_, err := analyzer.Analyze(context.Background(), imageData)
	if err == nil {
		t.Error("expected error from API, got nil")
	}
}

// TestGoogleVisionAnalyzer_Close tests Close method
func TestGoogleVisionAnalyzer_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCfg := cfgmocks.NewMockConfig(ctrl)
	mockClient := mocks.NewMockVisionClient(ctrl)

	// Expect Close to be called
	mockClient.EXPECT().
		Close().
		Return(nil)

	analyzer := NewGoogleVisionAnalyzerWithClient(mockClient, mockCfg)

	err := analyzer.Close()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestGoogleVisionAnalyzer_AnalyzeNoLabelsDetected tests when API returns empty labels
func TestGoogleVisionAnalyzer_AnalyzeNoLabelsDetected(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCfg := cfgmocks.NewMockConfig(ctrl)
	mockClient := mocks.NewMockVisionClient(ctrl)

	imageData := []byte("test image data")

	// Set up mock response with no labels
	expectedResponse := &visionpb.BatchAnnotateImagesResponse{
		Responses: []*visionpb.AnnotateImageResponse{
			{
				LabelAnnotations: nil,
			},
		},
	}

	// Expect BatchAnnotateImages to be called
	mockClient.EXPECT().
		BatchAnnotateImages(gomock.Any(), gomock.Any()).
		Return(expectedResponse, nil)

	analyzer := NewGoogleVisionAnalyzerWithClient(mockClient, mockCfg)

	_, err := analyzer.Analyze(context.Background(), imageData)
	if err == nil {
		t.Error("expected error when no labels detected, got nil")
	}
}

// TestGoogleVisionAnalyzer_AnalyzeEmptyResponses tests when API returns empty response list
func TestGoogleVisionAnalyzer_AnalyzeEmptyResponses(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCfg := cfgmocks.NewMockConfig(ctrl)
	mockClient := mocks.NewMockVisionClient(ctrl)

	imageData := []byte("test image data")

	// Set up mock response with empty responses
	expectedResponse := &visionpb.BatchAnnotateImagesResponse{
		Responses: []*visionpb.AnnotateImageResponse{},
	}

	// Expect BatchAnnotateImages to be called
	mockClient.EXPECT().
		BatchAnnotateImages(gomock.Any(), gomock.Any()).
		Return(expectedResponse, nil)

	analyzer := NewGoogleVisionAnalyzerWithClient(mockClient, mockCfg)

	_, err := analyzer.Analyze(context.Background(), imageData)
	if err == nil {
		t.Error("expected error when no responses returned, got nil")
	}
}

// TestGoogleVisionAnalyzer_AnalyzeWithInvalidTags tests label filtering
func TestGoogleVisionAnalyzer_AnalyzeWithInvalidTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCfg := cfgmocks.NewMockConfig(ctrl)
	mockClient := mocks.NewMockVisionClient(ctrl)

	imageData := []byte("test image data")

	// Set up mock response with invalid score (> 1.0)
	expectedResponse := &visionpb.BatchAnnotateImagesResponse{
		Responses: []*visionpb.AnnotateImageResponse{
			{
				LabelAnnotations: []*visionpb.EntityAnnotation{
					{
						Description: "label1",
						Score:       1.5, // Invalid score > 1.0
					},
				},
			},
		},
	}

	// Expect BatchAnnotateImages to be called
	mockClient.EXPECT().
		BatchAnnotateImages(gomock.Any(), gomock.Any()).
		Return(expectedResponse, nil)

	analyzer := NewGoogleVisionAnalyzerWithClient(mockClient, mockCfg)

	_, err := analyzer.Analyze(context.Background(), imageData)
	if err == nil {
		t.Error("expected error when all tags are invalid, got nil")
	}
}

// TestGoogleVisionAnalyzer_AnalyzeWithMixedTags tests mixed valid and invalid tags
func TestGoogleVisionAnalyzer_AnalyzeWithMixedTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCfg := cfgmocks.NewMockConfig(ctrl)
	mockClient := mocks.NewMockVisionClient(ctrl)

	imageData := []byte("test image data")

	// Set up mock response with mixed valid/invalid tags
	expectedResponse := &visionpb.BatchAnnotateImagesResponse{
		Responses: []*visionpb.AnnotateImageResponse{
			{
				LabelAnnotations: []*visionpb.EntityAnnotation{
					{
						Description: "label1",
						Score:       1.5, // Invalid - too high
					},
					{
						Description: "label2",
						Score:       0.85, // Valid
					},
					{
						Description: "label3",
						Score:       -0.1, // Invalid - negative
					},
				},
			},
		},
	}

	// Expect BatchAnnotateImages to be called
	mockClient.EXPECT().
		BatchAnnotateImages(gomock.Any(), gomock.Any()).
		Return(expectedResponse, nil)

	analyzer := NewGoogleVisionAnalyzerWithClient(mockClient, mockCfg)

	result, err := analyzer.Analyze(context.Background(), imageData)
	if err != nil {
		t.Errorf("expected no error for mixed tags, got %v", err)
	}

	if result == nil {
		t.Error("expected non-nil result")
	}

	// Should only have the valid tag
	if len(result.Tags) != 1 {
		t.Errorf("expected 1 valid tag, got %d", len(result.Tags))
	}
}

// TestGoogleVisionAnalyzer_CloseWithError tests Close method when client returns error
func TestGoogleVisionAnalyzer_CloseWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCfg := cfgmocks.NewMockConfig(ctrl)
	mockClient := mocks.NewMockVisionClient(ctrl)

	// Expect Close to be called and return error
	mockClient.EXPECT().
		Close().
		Return(errors.New("close error"))

	analyzer := NewGoogleVisionAnalyzerWithClient(mockClient, mockCfg)

	err := analyzer.Close()
	if err == nil {
		t.Error("expected error from Close, got nil")
	}
}
