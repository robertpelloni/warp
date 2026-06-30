package pty

import (
	"os/exec"
	"testing"
)

func TestNewPTY(t *testing.T) {
	shell := "echo"
	if _, err := exec.LookPath(shell); err != nil {
		t.Skip("echo not found in PATH")
	}

	cfg := Config{
		Shell: shell,
		Args:  []string{"hello"},
	}

	p, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer p.Close()

	if p == nil {
		t.Fatal("expected PTY, got nil")
	}

	cols, rows := p.Size()
	if cols != 120 || rows != 40 {
		t.Errorf("expected 120x40, got %dx%d", cols, rows)
	}

	err = p.Resize(80, 24)
	if err != nil {
		t.Errorf("unexpected error on resize: %v", err)
	}

	cols, rows = p.Size()
	if cols != 80 || rows != 24 {
		t.Errorf("expected 80x24, got %dx%d", cols, rows)
	}
}
