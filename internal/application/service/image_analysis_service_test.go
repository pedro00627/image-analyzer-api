package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pedro00627/image-analyzer-api/internal/application/usecase"
	"github.com/pedro00627/image-analyzer-api/internal/application/usecase/mocks"
	"github.com/pedro00627/image-analyzer-api/internal/domain/entity"
)

func TestNewImageAnalysisService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAnalyzeImageUseCaseInterface(ctrl)
	service := NewImageAnalysisService(mockUseCase)

	if service == nil {
		t.Error("NewImageAnalysisService() returned nil")
	}
}

func TestImageAnalysisService_AnalyzeImage_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAnalyzeImageUseCaseInterface(ctrl)
	service := NewImageAnalysisService(mockUseCase)

	expectedResult := &entity.AnalysisResult{
		Tags: []entity.Tag{
			{Label: "cat", Confidence: 0.92},
		},
	}

	mockUseCase.EXPECT().
		Execute(gomock.Any(), gomock.Any()).
		Return(expectedResult, nil)

	req := &AnalyzeImageRequest{
		Filename:  "cat.jpg",
		MimeType:  "image/jpeg",
		ImageData: []byte("fake image"),
	}

	result, err := service.AnalyzeImage(context.Background(), req)

	if err != nil {
		t.Errorf("AnalyzeImage() returned error: %v", err)
	}

	if result != expectedResult {
		t.Error("AnalyzeImage() did not return expected result")
	}
}

func TestImageAnalysisService_AnalyzeImage_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAnalyzeImageUseCaseInterface(ctrl)
	service := NewImageAnalysisService(mockUseCase)

	expectedError := errors.New("validation error")

	mockUseCase.EXPECT().
		Execute(gomock.Any(), gomock.Any()).
		Return(nil, expectedError)

	req := &AnalyzeImageRequest{
		Filename:  "invalid.txt",
		MimeType:  "text/plain",
		ImageData: []byte("not an image"),
	}

	result, err := service.AnalyzeImage(context.Background(), req)

	if err == nil {
		t.Error("AnalyzeImage() should return error")
	}

	if result != nil {
		t.Error("AnalyzeImage() should return nil result on error")
	}

	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}
}

func TestImageAnalysisService_AnalyzeImage_MapsRequestCorrectly(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAnalyzeImageUseCaseInterface(ctrl)
	service := NewImageAnalysisService(mockUseCase)

	expectedFilename := "dog.png"
	expectedMimeType := "image/png"
	expectedData := []byte("image data")

	mockUseCase.EXPECT().
		Execute(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, req usecase.AnalyzeImageRequest) (*entity.AnalysisResult, error) {
			if req.Filename != expectedFilename {
				t.Errorf("Expected filename %s, got %s", expectedFilename, req.Filename)
			}
			if req.MimeType != expectedMimeType {
				t.Errorf("Expected mime type %s, got %s", expectedMimeType, req.MimeType)
			}
			if string(req.ImageData) != string(expectedData) {
				t.Error("ImageData not mapped correctly")
			}
			return &entity.AnalysisResult{}, nil
		})

	req := &AnalyzeImageRequest{
		Filename:  expectedFilename,
		MimeType:  expectedMimeType,
		ImageData: expectedData,
	}

	service.AnalyzeImage(context.Background(), req)
}
