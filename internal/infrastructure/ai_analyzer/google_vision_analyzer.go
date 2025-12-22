package ai_analyzer

import (
	"context"
	"fmt"

	vision "cloud.google.com/go/vision/v2/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"github.com/pedro00627/image-analyzer-api/internal/domain/entity"
	"github.com/pedro00627/image-analyzer-api/internal/infrastructure/config"
)

//go:generate mockgen -source=./google_vision_analyzer.go -destination=mocks/mock_analyzer.go -package=mocks

type Analyzer interface {
	Analyze(ctx context.Context, imageData []byte) (*entity.AnalysisResult, error)
}

// GoogleVisionAnalyzer implements Analyzer using Google Cloud Vision API
type GoogleVisionAnalyzer struct {
	client VisionClient
	cfg    config.Config
}

var _ Analyzer = (*GoogleVisionAnalyzer)(nil)

// NewGoogleVisionAnalyzer creates a new GoogleVisionAnalyzer
// Credentials are loaded from GOOGLE_APPLICATION_CREDENTIALS environment variable
func NewGoogleVisionAnalyzer(ctx context.Context, cfg config.Config) (*GoogleVisionAnalyzer, error) {
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create vision client: %w", err)
	}

	return &GoogleVisionAnalyzer{
		client: client,
		cfg:    cfg,
	}, nil
}

// NewGoogleVisionAnalyzerWithClient creates a GoogleVisionAnalyzer with a custom client (for testing)
func NewGoogleVisionAnalyzerWithClient(client VisionClient, cfg config.Config) *GoogleVisionAnalyzer {
	return &GoogleVisionAnalyzer{
		client: client,
		cfg:    cfg,
	}
}

// Analyze analyzes an image using Google Cloud Vision API and returns detected labels
func (g *GoogleVisionAnalyzer) Analyze(ctx context.Context, imageData []byte) (*entity.AnalysisResult, error) {
	if imageData == nil || len(imageData) == 0 {
		return nil, fmt.Errorf("image data is empty")
	}

	// Create Image object from bytes
	image := &visionpb.Image{
		Content: imageData,
	}

	// Create label detection feature
	feature := &visionpb.Feature{
		Type: visionpb.Feature_LABEL_DETECTION,
	}

	// Create annotation request for single image
	request := &visionpb.AnnotateImageRequest{
		Image:    image,
		Features: []*visionpb.Feature{feature},
	}

	// Create batch request with single image
	batchRequest := &visionpb.BatchAnnotateImagesRequest{
		Requests: []*visionpb.AnnotateImageRequest{request},
	}

	// Make API call
	response, err := g.client.BatchAnnotateImages(ctx, batchRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to annotate image: %w", err)
	}

	// Extract labels from response
	tags := make([]entity.Tag, 0)
	if len(response.Responses) > 0 && response.Responses[0].LabelAnnotations != nil {
		for _, annotation := range response.Responses[0].LabelAnnotations {
			tag := entity.NewTag(annotation.Description, float64(annotation.Score))
			if tag.IsValid() {
				tags = append(tags, tag)
			}
		}
	}

	if len(tags) == 0 {
		return nil, fmt.Errorf("no labels detected in image")
	}

	// Create and return AnalysisResult
	result := entity.NewAnalysisResult(tags)
	return &result, nil
}

// Close closes the Vision API client
func (g *GoogleVisionAnalyzer) Close() error {
	if g.client != nil {
		return g.client.Close()
	}
	return nil
}
