// Package service contains domain service interfaces
package service

// ImageValidator defines the interface for image validation
// Methods return domain-level errors defined in internal/domain/error
// Size unit is bytes; type is MIME and/or filename check.
type ImageValidator interface {
	// ValidateType checks if the file type is allowed
	ValidateType(filename string, contentType string) error

	// ValidateSize checks if the file size is within limits
	ValidateSize(size int64) error

	// ValidateContent checks if the image content is valid and decodable
	ValidateContent(data []byte) error
}
