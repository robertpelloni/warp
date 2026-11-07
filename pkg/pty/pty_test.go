package pty

import (
	"runtime"
	"testing"
)

func TestPTYCreationAndFallback(t *testing.T) {
	sh := "sh"
	if runtime.GOOS == "windows" {
		sh = "cmd.exe"
	}
	p, err := New(Config{
		Shell: sh,
		Cols:  80,
		Rows:  24,
	})
	if err != nil {
		t.Fatalf("failed to create PTY: %v", err)
	}
	defer p.Close()

	if !p.Running() {
		t.Errorf("PTY expected to be running")
	}

	c, r := p.Size()
	if c != 80 || r != 24 {
		t.Errorf("expected size (80, 24), got (%d, %d)", c, r)
	}

	if err := p.Resize(100, 30); err != nil {
		t.Errorf("resize failed: %v", err)
	}
	c, r = p.Size()
	if c != 100 || r != 30 {
		t.Errorf("expected size (100, 30), got (%d, %d)", c, r)
	}
}
