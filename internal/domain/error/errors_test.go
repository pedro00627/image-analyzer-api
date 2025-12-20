package error

import "testing"

const testInvalidInputMessage = "Invalid input provided"

func TestDomainErrorError(t *testing.T) {
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
	err := NewDomainError("INVALID_INPUT", testInvalidInputMessage)

	if err.Code != "INVALID_INPUT" {
		t.Errorf("Code = %q, expected %q", err.Code, "INVALID_INPUT")
	}
	if err.Message != testInvalidInputMessage {
		t.Errorf("Message = %q, expected %q", err.Message, testInvalidInputMessage)
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
