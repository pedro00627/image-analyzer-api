package apierror

// ErrorCode represents API error codes
type ErrorCode string

const (
	InvalidRequest      ErrorCode = "invalid_request"
	InvalidFileType     ErrorCode = "invalid_file_type"
	FileTooLarge        ErrorCode = "file_too_large"
	InvalidImageContent ErrorCode = "invalid_image_content"
	AnalysisFailed      ErrorCode = "analysis_failed"
	InternalError       ErrorCode = "internal_error"
)
