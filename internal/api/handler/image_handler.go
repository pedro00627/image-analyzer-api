package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pedro00627/image-analyzer-api/internal/application/service"
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
	var req service.AnalyzeImageRequest

	// Parse multipart form
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing image file"})
		return
	}

	// Read file content
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	defer src.Close()

	// Read file data
	buf := make([]byte, file.Size)
	if _, err := src.Read(buf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file data"})
		return
	}

	req.Filename = file.Filename
	req.MimeType = file.Header.Get("Content-Type")
	req.ImageData = buf

	// Call service
	result, err := h.service.AnalyzeImage(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
