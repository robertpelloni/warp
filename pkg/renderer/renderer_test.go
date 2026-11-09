package renderer

import (
	"testing"
)

func TestNewRenderer(t *testing.T) {
	r := New(800, 600)
	if r == nil {
		t.Fatal("New returned nil")
	}
	if r.Atlas == nil {
		t.Fatal("Atlas not initialized")
	}
	r.Resize(1024, 768)
	if r.width != 1024 || r.height != 768 {
		t.Errorf("Resize failed: got %dx%d", r.width, r.height)
	}
}

func TestTextureAtlasAllocAndGet(t *testing.T) {
	atlas := NewTextureAtlas(128, 128)

	g, err := atlas.Alloc('A', 10, 16, 10)
	if err != nil {
		t.Fatalf("Alloc failed: %v", err)
	}
	if g.Rune != 'A' || g.Width != 10 || g.Height != 16 {
		t.Errorf("Allocated glyph mismatch: %+v", g)
	}

	cached, ok := atlas.Get('A')
	if !ok || cached != g {
		t.Errorf("Get('A') failed: got %+v, ok=%v", cached, ok)
	}

	if g.UVMaxX <= g.UVMinX || g.UVMaxY <= g.UVMinY {
		t.Errorf("Invalid UV coordinates: [%f, %f, %f, %f]", g.UVMinX, g.UVMinY, g.UVMaxX, g.UVMaxY)
	}
}

func TestTextureAtlasOverflow(t *testing.T) {
	atlas := NewTextureAtlas(32, 32)
	// Fill the small atlas
	_, err := atlas.Alloc('A', 20, 20, 20)
	if err != nil {
		t.Fatalf("First alloc failed: %v", err)
	}
	_, err = atlas.Alloc('B', 20, 20, 20)
	if err == nil {
		t.Error("Expected error when overflowing atlas bounds")
	}
}
