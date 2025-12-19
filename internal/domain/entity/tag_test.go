package entity

import (
	"testing"
)

func TestNewTag(t *testing.T) {
	tag := NewTag("Dog", 0.95)

	if tag.Label != "Dog" {
		t.Errorf("expected label 'Dog', got %s", tag.Label)
	}
	if tag.Confidence != 0.95 {
		t.Errorf("expected confidence 0.95, got %f", tag.Confidence)
	}
}

func TestTagIsValid(t *testing.T) {
	tests := []struct {
		name     string
		tag      Tag
		expected bool
	}{
		{
			name:     "valid tag",
			tag:      Tag{Label: "Cat", Confidence: 0.8},
			expected: true,
		},
		{
			name:     "empty label",
			tag:      Tag{Label: "", Confidence: 0.8},
			expected: false,
		},
		{
			name:     "negative confidence",
			tag:      Tag{Label: "Dog", Confidence: -0.1},
			expected: false,
		},
		{
			name:     "confidence over 1",
			tag:      Tag{Label: "Dog", Confidence: 1.1},
			expected: false,
		},
		{
			name:     "zero confidence",
			tag:      Tag{Label: "Dog", Confidence: 0.0},
			expected: true,
		},
		{
			name:     "confidence 1.0",
			tag:      Tag{Label: "Dog", Confidence: 1.0},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.tag.IsValid()
			if result != tt.expected {
				t.Errorf("IsValid() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
