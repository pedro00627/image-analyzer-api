package error

import "testing"

func TestDomainError_Error(t *testing.T) {
	err := &DomainError{
		Code:    "TEST_ERROR",
		Message: "This is a test error",
	}

	expected := "TEST_ERROR: This is a test error"
	if err.Error() != expected {
		t.Errorf("Error() = %q, expected %q", err.Error(), expected)
	}
}

func TestNewDomainError(t *testing.T) {
	err := NewDomainError("INVALID_INPUT", "Invalid input provided")

	if err.Code != "INVALID_INPUT" {
		t.Errorf("Code = %q, expected %q", err.Code, "INVALID_INPUT")
	}
	if err.Message != "Invalid input provided" {
		t.Errorf("Message = %q, expected %q", err.Message, "Invalid input provided")
	}
}

func TestPredefinedErrors(t *testing.T) {
	if ErrInvalidFileType.Code != "INVALID_FILE_TYPE" {
		t.Errorf("ErrInvalidFileType.Code = %q", ErrInvalidFileType.Code)
	}
	if ErrFileTooLarge.Code != "FILE_TOO_LARGE" {
		t.Errorf("ErrFileTooLarge.Code = %q", ErrFileTooLarge.Code)
	}
	// Add more as needed
}
