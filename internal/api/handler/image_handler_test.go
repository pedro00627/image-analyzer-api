package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/pedro00627/image-analyzer-api/internal/application/service/mocks"
	"github.com/pedro00627/image-analyzer-api/internal/domain/entity"
	domainError "github.com/pedro00627/image-analyzer-api/internal/domain/error"
)

func TestNewImageHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockImageAnalysisService(ctrl)
	handler := NewImageHandler(mockService)

	if handler == nil {
		t.Fatal("NewImageHandler() returned nil")
	}

	if handler.service != mockService {
		t.Error("NewImageHandler() did not set service correctly")
	}
}

func TestImageHandler_AnalyzeImage_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockImageAnalysisService(ctrl)
	handler := NewImageHandler(mockService)

	expectedResult := &entity.AnalysisResult{
		Tags: []entity.Tag{
			{Label: "dog", Confidence: 0.95},
		},
	}

	mockService.EXPECT().
		AnalyzeImage(gomock.Any(), gomock.Any()).
		Return(expectedResult, nil)

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", "test.jpg")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	_, err = part.Write([]byte("fake image data"))
	if err != nil {
		t.Fatalf("Failed to write to form file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}

	// Create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("POST", "/api/analyze", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.Request = req

	// Execute
	handler.AnalyzeImage(c)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || !success {
		t.Error("Expected success: true")
	}
}

func TestImageHandler_AnalyzeImage_MissingFile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockImageAnalysisService(ctrl)
	handler := NewImageHandler(mockService)

	// Create request without file
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("POST", "/api/analyze", nil)
	c.Request = req

	// Execute
	handler.AnalyzeImage(c)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || success {
		t.Error("Expected success: false")
	}
}

func TestImageHandler_AnalyzeImage_DomainError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockImageAnalysisService(ctrl)
	handler := NewImageHandler(mockService)

	mockService.EXPECT().
		AnalyzeImage(gomock.Any(), gomock.Any()).
		Return(nil, domainError.ErrInvalidFileType)

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", "test.jpg")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	_, err = part.Write([]byte("fake image data"))
	if err != nil {
		t.Fatalf("Failed to write to form file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}

	// Create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("POST", "/api/analyze", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.Request = req

	// Execute
	handler.AnalyzeImage(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || success {
		t.Error("Expected success: false")
	}

	if errorObj, ok := response["error"].(map[string]interface{}); ok {
		if code := errorObj["code"].(string); code != domainError.CodeInvalidFileType {
			t.Errorf("Expected error code %s, got %s", domainError.CodeInvalidFileType, code)
		}
	} else {
		t.Error("Expected error object in response")
	}
}

func TestImageHandler_AnalyzeImage_UnknownError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockImageAnalysisService(ctrl)
	handler := NewImageHandler(mockService)

	mockService.EXPECT().
		AnalyzeImage(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("unknown error"))

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", "test.jpg")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	_, err = part.Write([]byte("fake image data"))
	if err != nil {
		t.Fatalf("Failed to write to form file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}

	// Create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("POST", "/api/analyze", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.Request = req

	// Execute
	handler.AnalyzeImage(c)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || success {
		t.Error("Expected success: false")
	}

	if errorObj, ok := response["error"].(map[string]interface{}); ok {
		if code, ok := errorObj["code"].(string); ok {
			if code != "internal_error" {
				t.Errorf("Expected error code internal_error, got %s", code)
			}
		} else {
			t.Error("Expected error code string in error object")
		}
	} else {
		t.Error("Expected error object in response")
	}
}

func TestImageHandler_parseRequest_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockImageAnalysisService(ctrl)
	handler := NewImageHandler(mockService)

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", "test.jpg")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	testData := []byte("fake image data")
	_, err = part.Write(testData)
	if err != nil {
		t.Fatalf("Failed to write to form file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}

	// Create request
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	req := httptest.NewRequest("POST", "/api/analyze", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.Request = req

	// Execute
	result, err := handler.parseRequest(c)

	// Assert
	if err != nil {
		t.Errorf("parseRequest() returned error: %v", err)
	}

	if result == nil {
		t.Fatal("parseRequest() returned nil result")
	}

	if result.Filename != "test.jpg" {
		t.Errorf("Expected filename test.jpg, got %s", result.Filename)
	}

	if !bytes.Equal(result.ImageData, testData) {
		t.Error("ImageData does not match")
	}
}

func TestImageHandler_parseRequest_MissingFile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockImageAnalysisService(ctrl)
	handler := NewImageHandler(mockService)

	// Create request without file
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	req := httptest.NewRequest("POST", "/api/analyze", nil)
	c.Request = req

	// Execute
	result, err := handler.parseRequest(c)

	// Assert
	if err == nil {
		t.Error("parseRequest() should return error for missing file")
	}

	if result != nil {
		t.Error("parseRequest() should return nil result on error")
	}
}

func TestImageHandler_handleError_DomainError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockImageAnalysisService(ctrl)
	handler := NewImageHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.handleError(c, domainError.ErrFileTooLarge)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if errorObj, ok := response["error"].(map[string]interface{}); ok {
		if code, ok := errorObj["code"].(string); ok {
			if code != domainError.CodeFileTooLarge {
				t.Errorf("Expected error code %s, got %s", domainError.CodeFileTooLarge, code)
			}
		} else {
			t.Error("Expected error code string in error object")
		}
	} else {
		t.Error("Expected error object in response")
	}
}

func TestImageHandler_handleError_UnknownError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockImageAnalysisService(ctrl)
	handler := NewImageHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.handleError(c, errors.New("test error"))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if errorObj, ok := response["error"].(map[string]interface{}); ok {
		if code, ok := errorObj["code"].(string); ok {
			if code != "internal_error" {
				t.Errorf("Expected error code internal_error, got %s", code)
			}
		} else {
			t.Error("Expected error code string in error object")
		}
	} else {
		t.Error("Expected error object in response")
	}
}
