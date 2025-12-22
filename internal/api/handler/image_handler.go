package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pedro00627/image-analyzer-api/internal/application/service"
	domainError "github.com/pedro00627/image-analyzer-api/internal/domain/error"
)

// ImageHandler handles image analysis requests
type ImageHandler struct {
	service service.ImageAnalysisService
}

// NewImageHandler creates a new ImageHandler
func NewImageHandler(svc service.ImageAnalysisService) *ImageHandler {
	return &ImageHandler{
		service: svc,
	}
}

// AnalyzeImage handles POST /api/analyze requests
func (h *ImageHandler) AnalyzeImage(c *gin.Context) {
	// Parse request
	req, err := h.parseRequest(c)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Call service
	result, err := h.service.AnalyzeImage(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// parseRequest extracts and validates the image from the request
func (h *ImageHandler) parseRequest(c *gin.Context) (*service.AnalyzeImageRequest, error) {
	file, err := c.FormFile("image")
	if err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	return &service.AnalyzeImageRequest{
		Filename:  file.Filename,
		MimeType:  file.Header.Get("Content-Type"),
		ImageData: data,
	}, nil
}

// handleError maps domain errors to HTTP responses
func (h *ImageHandler) handleError(c *gin.Context, err error) {
	if domainErr, ok := err.(*domainError.DomainError); ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    domainErr.Code,
				"message": domainErr.Message,
			},
		})
		return
	}

	// Unknown error
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "internal_error",
			"message": err.Error(),
		},
	})
}
