package entity

import (
	"testing"
	"time"
)

func TestNewAnalysisResult(t *testing.T) {
	tags := []Tag{
		NewTag("Dog", 0.95),
		NewTag("Cat", 0.8),
	}
	result := NewAnalysisResult(tags)

	if len(result.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(result.Tags))
	}
	if result.Tags[0].Label != "Dog" {
		t.Errorf("expected first tag 'Dog', got %s", result.Tags[0].Label)
	}
	if result.AnalyzedAt.IsZero() {
		t.Error("AnalyzedAt should not be zero")
	}
}

func TestAnalysisResultIsValid(t *testing.T) {
	tests := []struct {
		name     string
		result   AnalysisResult
		expected bool
	}{
		{
			name: "valid result",
			result: AnalysisResult{
				Tags: []Tag{
					NewTag("Dog", 0.95),
					NewTag("Cat", 0.8),
				},
				AnalyzedAt: time.Now(),
			},
			expected: true,
		},
		{
			name: "empty tags",
			result: AnalysisResult{
				Tags:       []Tag{},
				AnalyzedAt: time.Now(),
			},
			expected: false,
		},
		{
			name: "invalid tag",
			result: AnalysisResult{
				Tags: []Tag{
					NewTag("", 0.95), // invalid label
				},
				AnalyzedAt: time.Now(),
			},
			expected: false,
		},
		{
			name: "mixed valid and invalid",
			result: AnalysisResult{
				Tags: []Tag{
					NewTag("Dog", 0.95),
					NewTag("", 0.8), // invalid
				},
				AnalyzedAt: time.Now(),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.result.IsValid() != tt.expected {
				t.Errorf("IsValid() = %v, expected %v", tt.result.IsValid(), tt.expected)
			}
		})
	}
}
