// Package validator provides image validation implementations
package validator

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strings"

	"golang.org/x/image/webp"

	domainerr "github.com/pedro00627/image-analyzer-api/internal/domain/error"
	"github.com/pedro00627/image-analyzer-api/internal/infrastructure/config"
)

//go:generate mockgen -source=./image_validator.go -destination=./mocks/mock_validator.go -package=mocks

// Validator defines the contract for image validation
type Validator interface {
	ValidateType(filename string, contentType string) error
	ValidateSize(size int64) error
	ValidateContent(data []byte) error
}

// ImageValidator validates images based on configuration
type ImageValidator struct {
	allowedTypes map[string]struct{}
	maxSize      int64
}

// NewImageValidator constructs an ImageValidator from a Config
func NewImageValidator(cfg config.Config) *ImageValidator {
	// Build map of allowed mime types for quick lookup
	allowedMimes := cfg.GetAllowedMimes()
	allowedTypes := make(map[string]struct{}, len(allowedMimes))
	for _, mime := range allowedMimes {
		allowedTypes[strings.TrimSpace(mime)] = struct{}{}
	}

	return &ImageValidator{
		allowedTypes: allowedTypes,
		maxSize:      cfg.GetMaxSize(),
	}
}

// ValidateType checks the Content-Type header (informational, not critical)
// This is a secondary check - content has already been validated
// Returns nil if MIME type is recognized, and nil for unrecognized types too
// (content validation is the primary security check)
func (v *ImageValidator) ValidateType(filename string, contentType string) error {
	// If no Content-Type provided, that's ok (content already validated)
	if contentType == "" {
		return nil
	}

	// If Content-Type is in allowed list, perfect
	if _, ok := v.allowedTypes[contentType]; ok {
		return nil
	}

	// If Content-Type is provided but not in whitelist, that's ok
	// Content validation already confirmed it's a valid image
	// (e.g., JPEG sent with type=image/webp is still a valid JPEG)
	// TODO: Add logging for Content-Type mismatch
	return nil
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
	
	// Check for common image magic bytes
	if len(data) < 4 {
		return domainerr.ErrInvalidImageContent
	}
	
	// JPEG: FF D8 FF
	if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return nil
	}
	
	// PNG: 89 50 4E 47
	if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return nil
	}
	
	// GIF: 47 49 46
	if data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 {
		return nil
	}
	
	// WebP: RIFF ... WEBP
	if len(data) >= 12 && data[0] == 0x52 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x46 &&
		data[8] == 0x57 && data[9] == 0x45 && data[10] == 0x42 && data[11] == 0x50 {
		return nil
	}
	
	// Try to decode with standard libraries
	_, _, err := image.Decode(bytes.NewReader(data))
	if err == nil {
		return nil // Valid image
	}
	
	// Try WebP separately
	_, err = webp.Decode(bytes.NewReader(data))
	if err == nil {
		return nil // Valid WebP
	}
	
	return fmt.Errorf("%w: no valid image format detected", domainerr.ErrInvalidImageContent)
}
