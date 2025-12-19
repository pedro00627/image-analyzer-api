// Package entity contains domain entities and value objects
package entity

import "time"

// AnalysisResult represents the result of analyzing an image
type AnalysisResult struct {
	Tags      []Tag     `json:"tags"`
	AnalyzedAt time.Time `json:"analyzed_at"`
}

// NewAnalysisResult creates a new AnalysisResult
func NewAnalysisResult(tags []Tag) AnalysisResult {
	return AnalysisResult{
		Tags:      tags,
		AnalyzedAt: time.Now(),
	}
}

// IsValid checks if the result has valid data
func (r AnalysisResult) IsValid() bool {
	if len(r.Tags) == 0 {
		return false
	}
	for _, tag := range r.Tags {
		if !tag.IsValid() {
			return false
		}
	}
	return true
}