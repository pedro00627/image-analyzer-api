package config

import "testing"

func TestEnvReader_ReturnsEnvValue(t *testing.T) {
	r := envReader{}
	t.Setenv("CFG_KEY", "value")

	got := r.Get("CFG_KEY", "default")
	if got != "value" {
		t.Fatalf("Get() = %q, want %q", got, "value")
	}
}

func TestEnvReader_ReturnsDefaultWhenUnset(t *testing.T) {
	r := envReader{}
	t.Setenv("CFG_MISSING", "")

	got := r.Get("CFG_MISSING", "default")
	if got != "default" {
		t.Fatalf("Get() = %q, want %q", got, "default")
	}
}
