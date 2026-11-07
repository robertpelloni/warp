package app

import (
	"testing"
)

func TestAppConfig(t *testing.T) {
	cfg := Config{
		Shell:        "bash",
		TerminalOnly: true,
		Theme:        "Dracula",
	}
	a := New(cfg)
	if a == nil {
		t.Fatal("New returned nil")
	}
	if a.shell != "bash" {
		t.Errorf("shell mismatch: got %q, want 'bash'", a.shell)
	}
	if a.warpTheme == nil || a.warpTheme.Name() != "Dracula" {
		t.Errorf("theme mismatch: got %v", a.warpTheme)
	}
	a.Quit()
}
