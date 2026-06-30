package renderer

import (
	"testing"
)

func TestNewRenderer(t *testing.T) {
	r := New(800, 600)
	if r == nil {
		t.Fatal("expected Renderer, got nil")
	}

	if r.width != 800 || r.height != 600 {
		t.Errorf("expected 800x600, got %dx%d", r.width, r.height)
	}

	r.Resize(1024, 768)
	if r.width != 1024 || r.height != 768 {
		t.Errorf("expected 1024x768, got %dx%d", r.width, r.height)
	}
}
