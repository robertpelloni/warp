package editor

import (
	"testing"
)

func TestNew(t *testing.T) {
	e := New()
	if e == nil {
		t.Fatal("New() returned nil")
	}
}

func TestOpenBuffer(t *testing.T) {
	e := New()
	b := e.Open("test.go")
	if b == nil {
		t.Fatal("Open() returned nil")
	}
	if b.Path != "test.go" {
		t.Errorf("Path = %q, want 'test.go'", b.Path)
	}
}

func TestOpenSameFile(t *testing.T) {
	e := New()
	b1 := e.Open("test.go")
	b1.Content = "hello"
	b2 := e.Open("test.go")
	if b2.Content != "hello" {
		t.Error("Opening same file should return existing buffer")
	}
}

func TestActiveBuffer(t *testing.T) {
	e := New()
	if e.Active() != nil {
		t.Error("Active() should be nil when no buffers are open")
	}
	e.Open("test.go")
	if e.Active() == nil {
		t.Error("Active() should not be nil after opening a buffer")
	}
}

func TestSetContent(t *testing.T) {
	e := New()
	e.Open("test.go")
	e.SetContent("test.go", "package main")
	b := e.Active()
	if b.Content != "package main" {
		t.Errorf("Content = %q, want 'package main'", b.Content)
	}
	if !b.Dirty {
		t.Error("Buffer should be dirty after SetContent")
	}
}

func TestSave(t *testing.T) {
	e := New()
	e.Open("test.go")
	e.SetContent("test.go", "package main")
	e.Save("test.go")
	b := e.Active()
	if b.Dirty {
		t.Error("Buffer should not be dirty after Save")
	}
}

func TestBuffers(t *testing.T) {
	e := New()
	e.Open("a.go")
	e.Open("b.go")
	buffers := e.Buffers()
	if len(buffers) != 2 {
		t.Errorf("Buffers() len = %d, want 2", len(buffers))
	}
}

func TestClose(t *testing.T) {
	e := New()
	e.Open("test.go")
	e.Close()
	if len(e.Buffers()) != 0 {
		t.Error("Buffers should be empty after Close")
	}
}

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"main.go", "go"},
		{"lib.rs", "rust"},
		{"script.py", "python"},
		{"app.js", "javascript"},
		{"app.ts", "javascript"},
		{"Main.java", "java"},
		{"style.css", "css"},
		{"page.html", "html"},
		{"config.json", "json"},
		{"data.yaml", "yaml"},
		{"data.yml", "yaml"},
		{"Cargo.toml", "toml"},
		{"run.sh", "bash"},
		{"README.md", "markdown"},
		{"query.sql", "sql"},
		{"header.h", "c"},
		{"main.c", "c"},
		{"app.cpp", "cpp"},
		{"Makefile", "text"},
	}

	for _, tt := range tests {
		result := DetectLanguage(tt.path)
		if result != tt.expected {
			t.Errorf("DetectLanguage(%q) = %q, want %q", tt.path, result, tt.expected)
		}
	}
}

func TestLineCount(t *testing.T) {
	b := &Buffer{Content: "line1\nline2\nline3"}
	if b.LineCount() != 3 {
		t.Errorf("LineCount() = %d, want 3", b.LineCount())
	}

	b2 := &Buffer{Content: ""}
	if b2.LineCount() != 1 {
		t.Errorf("LineCount() for empty = %d, want 1", b2.LineCount())
	}
}
