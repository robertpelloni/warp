package renderer

import (
	"fmt"
	"sync"
)

// Glyph represents a cached rasterized glyph in the texture atlas.
type Glyph struct {
	Rune     rune
	X        int
	Y        int
	Width    int
	Height   int
	AdvanceX int
	UVMinX   float32
	UVMinY   float32
	UVMaxX   float32
	UVMaxY   float32
}

// TextureAtlas packs glyphs into a 2D texture map.
type TextureAtlas struct {
	mu      sync.RWMutex
	Width   int
	Height  int
	CursorX int
	CursorY int
	RowMaxH int
	Glyphs  map[rune]*Glyph
}

// NewTextureAtlas creates a texture atlas with the given dimensions.
func NewTextureAtlas(w, h int) *TextureAtlas {
	if w <= 0 {
		w = 1024
	}
	if h <= 0 {
		h = 1024
	}
	return &TextureAtlas{
		Width:  w,
		Height: h,
		Glyphs: make(map[rune]*Glyph),
	}
}

// Alloc reserves space in the texture atlas for a glyph of given dimensions.
func (a *TextureAtlas) Alloc(r rune, gw, gh, advanceX int) (*Glyph, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if g, ok := a.Glyphs[r]; ok {
		return g, nil
	}

	if a.CursorX+gw > a.Width {
		a.CursorX = 0
		a.CursorY += a.RowMaxH
		a.RowMaxH = 0
	}

	if a.CursorY+gh > a.Height {
		return nil, fmt.Errorf("texture atlas out of space: %dx%d full", a.Width, a.Height)
	}

	glyph := &Glyph{
		Rune:     r,
		X:        a.CursorX,
		Y:        a.CursorY,
		Width:    gw,
		Height:   gh,
		AdvanceX: advanceX,
		UVMinX:   float32(a.CursorX) / float32(a.Width),
		UVMinY:   float32(a.CursorY) / float32(a.Height),
		UVMaxX:   float32(a.CursorX+gw) / float32(a.Width),
		UVMaxY:   float32(a.CursorY+gh) / float32(a.Height),
	}

	a.Glyphs[r] = glyph
	a.CursorX += gw
	if gh > a.RowMaxH {
		a.RowMaxH = gh
	}

	return glyph, nil
}

// Get retrieves a cached glyph.
func (a *TextureAtlas) Get(r rune) (*Glyph, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	g, ok := a.Glyphs[r]
	return g, ok
}

// Renderer manages GPU/backend rendering surfaces and font texture atlases.
type Renderer struct {
	width  int
	height int
	Atlas  *TextureAtlas
}

// New creates a new renderer with an initialized texture atlas.
func New(width, height int) *Renderer {
	return &Renderer{
		width:  width,
		height: height,
		Atlas:  NewTextureAtlas(1024, 1024),
	}
}

// Resize updates the renderer dimensions.
func (r *Renderer) Resize(w, h int) {
	r.width = w
	r.height = h
}
