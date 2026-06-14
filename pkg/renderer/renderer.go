package renderer

// Renderer is a placeholder for GPU-accelerated rendering.
// The actual rendering is handled by the Win32 GUI in pkg/app.
type Renderer struct {
	width  int
	height int
}

// New creates a new renderer.
func New(width, height int) *Renderer {
	return &Renderer{width: width, height: height}
}

// Resize updates the renderer dimensions.
func (r *Renderer) Resize(w, h int) {
	r.width = w
	r.height = h
}
