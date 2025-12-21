// Package validator provides image validation implementations
package validator

import (
	"bytes"
	"fmt"
	"image"
	"strings"

	domainerr "github.com/pedro00627/image-analyzer-api/internal/domain/error"
)

// ValidationConfig holds the configuration for image validation
type ValidationConfig struct {
	AllowedMimes []string // e.g., ["image/png", "image/jpeg"]
	MaxSize      int64    // maximum file size in bytes
}

// ImageValidator validates images based on configuration
type ImageValidator struct {
	allowedTypes map[string]struct{}
	maxSize      int64
}

// NewImageValidator constructs an ImageValidator from a ValidationConfig
func NewImageValidator(cfg ValidationConfig) *ImageValidator {
	// Build map of allowed mime types for quick lookup
	allowedTypes := make(map[string]struct{}, len(cfg.AllowedMimes))
	for _, mime := range cfg.AllowedMimes {
		allowedTypes[strings.TrimSpace(mime)] = struct{}{}
	}

	return &ImageValidator{
		allowedTypes: allowedTypes,
		maxSize:      cfg.MaxSize,
	}
}

// ValidateType checks if the file type (by MIME) is allowed
// Requires a recognized contentType; otherwise returns ErrInvalidFileType
func (v *ImageValidator) ValidateType(filename string, contentType string) error {
	// First: try MIME type if provided
	if contentType != "" {
		if _, ok := v.allowedTypes[contentType]; ok {
			return nil
		}
	}

	return domainerr.ErrInvalidFileType
}

// ValidateSize checks if the file size is within limits
func (v *ImageValidator) ValidateSize(size int64) error {
	if size > v.maxSize {
		return domainerr.ErrFileTooLarge
	}
	return nil
}

// ValidateContent decodes the image bytes to ensure it's a valid image
func (v *ImageValidator) ValidateContent(data []byte) error {
	if len(data) == 0 {
		return domainerr.ErrInvalidImageContent
	}
	_, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("%w: %v", domainerr.ErrInvalidImageContent, err)
	}
	return nil
}
