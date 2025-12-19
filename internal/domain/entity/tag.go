package entity

// Tag represents a label identified in an image with confidence score
type Tag struct {
	Label      string  `json:"label"`
	Confidence float64 `json:"confidence"`
}

func NewTag(label string, confidence float64) Tag {
	return Tag{
		Label:      label,
		Confidence: confidence,
	}
}

// IsValid checks if the tag has valid name and confidence score between 0 and 1
func (t Tag) IsValid() bool {
	return t.Label != "" && t.Confidence >= 0.0 && t.Confidence <= 1.0
}
