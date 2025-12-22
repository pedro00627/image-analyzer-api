package ai_analyzer

import (
	"context"

	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"github.com/googleapis/gax-go/v2"
)

//go:generate mockgen -source=./vision_client.go -destination=mocks/mock_vision_client.go -package=mocks

// VisionClient is an interface for Google Cloud Vision API client
type VisionClient interface {
	BatchAnnotateImages(context.Context, *visionpb.BatchAnnotateImagesRequest, ...gax.CallOption) (*visionpb.BatchAnnotateImagesResponse, error)
	Close() error
}
