package validator

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

// mockValidatorConfig implements config.Config for testing
type mockValidatorConfig struct {
	allowedMimes []string
	maxSize      int64
}

func (m *mockValidatorConfig) GetProvider() string {
	return "test"
}

func (m *mockValidatorConfig) GetSecret() string {
	return "test-secret"
}

func (m *mockValidatorConfig) GetAllowedMimes() []string {
	return m.allowedMimes
}

func (m *mockValidatorConfig) GetMaxSize() int64 {
	return m.maxSize
}

// helper to create a small valid PNG
func makePNGBytes(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			img.Set(x, y, color.RGBA{R: 255, A: 255})
		}
	}
	var buf []byte
	b := newBuffer(&buf)
	if err := png.Encode(b, img); err != nil {
		t.Fatalf("failed to encode png: %v", err)
	}
	return buf
}

// minimal writer to []byte without importing bytes
type byteBuffer struct{ data *[]byte }

func newBuffer(dst *[]byte) *byteBuffer { return &byteBuffer{data: dst} }

func (b *byteBuffer) Write(p []byte) (int, error) {
	*b.data = append(*b.data, p...)
	return len(p), nil
}

func TestValidateType(t *testing.T) {
	cfg := &mockValidatorConfig{
		allowedMimes: []string{"image/png", "image/jpeg"},
		maxSize:      10_000_000,
	}
	v := NewImageValidator(cfg)

	tests := []struct {
		name        string
		filename    string
		contentType string
		wantErr     bool
	}{
		{"png by mime", "file.bin", "image/png", false},
		{"jpeg by mime", "file.bin", "image/jpeg", false},
		{"missing mime", "photo.png", "", true},
		{"unsupported mime", "file.bin", "application/octet-stream", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateType(tt.filename, tt.contentType)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateSize(t *testing.T) {
	cfg := &mockValidatorConfig{
		allowedMimes: []string{"image/png"},
		maxSize:      5,
	}
	v := NewImageValidator(cfg)

	if err := v.ValidateSize(5); err != nil {
		t.Fatalf("ValidateSize() unexpected error: %v", err)
	}
	if err := v.ValidateSize(6); err == nil {
		t.Fatalf("ValidateSize() expected error for size over limit")
	}
}

func TestValidateContent(t *testing.T) {
	cfg := &mockValidatorConfig{
		allowedMimes: []string{"image/png"},
		maxSize:      10_000_000,
	}
	v := NewImageValidator(cfg)
	good := makePNGBytes(t)
	if err := v.ValidateContent(good); err != nil {
		t.Fatalf("ValidateContent() unexpected error: %v", err)
	}

	bad := []byte{0x00, 0x01, 0x02}
	if err := v.ValidateContent(bad); err == nil {
		t.Fatalf("ValidateContent() expected error for invalid image")
	}

	empty := []byte{}
	if err := v.ValidateContent(empty); err == nil {
		t.Fatalf("ValidateContent() expected error for empty data")
	}
}

// optional: ensure env does not leak (not used here but placeholder)
func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
