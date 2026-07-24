package editor

import (
	"fmt"
	"strings"
	"sync"
)

// Editor provides the Warp IDE editing capabilities.
type Editor struct {
	mu      sync.RWMutex
	buffers map[string]*Buffer
	active  string
}

// Buffer represents an open file.
type Buffer struct {
	Path     string
	Content  string
	Dirty    bool
	Cursor   CursorPos
	Language string
}

// CursorPos represents a cursor position.
type CursorPos struct {
	Line int
	Col  int
}

// New creates a new editor.
func New() *Editor {
	return &Editor{
		buffers: make(map[string]*Buffer),
	}
}

// Open opens a file buffer.
func (e *Editor) Open(path string) *Buffer {
	e.mu.Lock()
	defer e.mu.Unlock()

	if b, ok := e.buffers[path]; ok {
		e.active = path
		return b
	}

	b := &Buffer{
		Path:   path,
		Cursor: CursorPos{Line: 1, Col: 1},
	}
	e.buffers[path] = b
	e.active = path
	return b
}

// Active returns the active buffer.
func (e *Editor) Active() *Buffer {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.buffers[e.active]
}

// SetContent sets the content of a buffer.
func (e *Editor) SetContent(path, content string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if b, ok := e.buffers[path]; ok {
		b.Content = content
		b.Dirty = true
	}
}

// Save marks a buffer as saved.
func (e *Editor) Save(path string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if b, ok := e.buffers[path]; ok {
		b.Dirty = false
	}
}

// Close shuts down the editor.
func (e *Editor) Close() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.buffers = make(map[string]*Buffer)
}

// Buffers returns all open buffers.
func (e *Editor) Buffers() []*Buffer {
	e.mu.RLock()
	defer e.mu.RUnlock()
	result := make([]*Buffer, 0, len(e.buffers))
	for _, b := range e.buffers {
		result = append(result, b)
	}
	return result
}

// DetectLanguage returns a language identifier based on file extension.
func DetectLanguage(path string) string {
	parts := strings.Split(path, ".")
	if len(parts) < 2 {
		return "text"
	}
	ext := strings.ToLower(parts[len(parts)-1])
	switch ext {
	case "go":
		return "go"
	case "rs":
		return "rust"
	case "py":
		return "python"
	case "js", "ts":
		return "javascript"
	case "java":
		return "java"
	case "c", "h":
		return "c"
	case "cpp", "cc", "cxx":
		return "cpp"
	case "md":
		return "markdown"
	case "json":
		return "json"
	case "yaml", "yml":
		return "yaml"
	case "toml":
		return "toml"
	case "sh", "bash":
		return "bash"
	case "sql":
		return "sql"
	case "html":
		return "html"
	case "css":
		return "css"
	default:
		return "text"
	}
}

// LineCount returns the number of lines in a buffer.
func (b *Buffer) LineCount() int {
	return strings.Count(b.Content, "\n") + 1
}

// String returns a summary of the buffer.
func (b *Buffer) String() string {
	return fmt.Sprintf("%s (%d lines, dirty=%v)", b.Path, b.LineCount(), b.Dirty)
}
