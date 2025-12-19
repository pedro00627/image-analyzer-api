// Package error contains domain-specific errors
package error

// Error codes as constants
const (
	CodeInvalidFileType      = "INVALID_FILE_TYPE"
	CodeFileTooLarge         = "FILE_TOO_LARGE"
	CodeInvalidImageContent  = "INVALID_IMAGE_CONTENT"
	CodeAIServiceUnavailable = "AI_SERVICE_UNAVAILABLE"
)

// Error messages as constants
const (
	MsgInvalidFileType      = "Only JPEG, PNG, GIF, and WebP images are allowed"
	MsgFileTooLarge         = "File size exceeds 10MB limit"
	MsgInvalidImageContent  = "The uploaded file is not a valid image"
	MsgAIServiceUnavailable = "AI analysis service is currently unavailable"
)

// DomainError represents a domain-specific error with a code and message
type DomainError struct {
	Code    string
	Message string
}

// Error implements the error interface
func (e *DomainError) Error() string {
	return e.Code + ": " + e.Message
}

// NewDomainError creates a new DomainError
func NewDomainError(code, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}

// Predefined domain errors
var (
	ErrInvalidFileType = &DomainError{
		Code:    CodeInvalidFileType,
		Message: MsgInvalidFileType,
	}
	ErrFileTooLarge = &DomainError{
		Code:    CodeFileTooLarge,
		Message: MsgFileTooLarge,
	}
	ErrInvalidImageContent = &DomainError{
		Code:    CodeInvalidImageContent,
		Message: MsgInvalidImageContent,
	}
	ErrAIServiceUnavailable = &DomainError{
		Code:    CodeAIServiceUnavailable,
		Message: MsgAIServiceUnavailable,
	}
)
