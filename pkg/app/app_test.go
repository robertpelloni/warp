package app

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/robertpelloni/warp/pkg/command"
	"github.com/robertpelloni/warp/pkg/editor"
	"github.com/robertpelloni/warp/pkg/session"
	"github.com/robertpelloni/warp/pkg/shadowpilot"
)

func getTestShell() string {
	if runtime.GOOS == "windows" {
		return "cmd.exe"
	}
	return "sh"
}

func TestAppLifecycleIntegration(t *testing.T) {
	// Create required components
	sessMgr := session.NewManager()
	cmdEngine := command.NewEngine()
	editEngine := editor.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize shadow pilot
	var anomalyDetected bool
	pilot := shadowpilot.New(func(msg string) {
		_ = anomalyDetected
	})
	pilot.Start(ctx)

	// Create a new session
	sess, err := sessMgr.Create("Test Session", getTestShell(), 80, 24, nil)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}
	if sessMgr.Count() != 1 {
		t.Errorf("Expected 1 session, got %d", sessMgr.Count())
	}

	// Start terminal
	err = sess.Terminal.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start terminal: %v", err)
	}

	// Execute a command
	cmd := "hello"
	expanded := cmdEngine.Expand(cmd)
	if expanded != cmd {
		t.Errorf("Expected cmd to not expand, got %s", expanded)
	}

	cmdType := cmdEngine.Classify(cmd)
	if cmdType != command.CmdShell {
		t.Errorf("Expected CmdShell type, got %v", cmdType)
	}

	// Send command to terminal
	err = sess.Terminal.SendCommand(cmd)
	if err != nil {
		t.Errorf("Failed to send command to terminal: %v", err)
	}

	// Wait briefly to simulate running process
	time.Sleep(100 * time.Millisecond)

	// Check output buffer logic (it won't be filled due to dummy pipe, but should not crash)
	out := sess.Terminal.GetOutput()
	if len(out) > 0 {
		t.Logf("Terminal output: %s", out)
	}

	// Open an editor buffer
	buf := editEngine.Open("test.go")
	buf.Language = editor.DetectLanguage(buf.Path)
	if buf == nil {
		t.Fatal("Failed to open buffer")
	}

	// Validate language detection
	if buf.Language != "go" {
		t.Errorf("Expected 'go' language, got '%s'", buf.Language)
	}

	// Terminate
	sess.Terminal.Stop()
	sessMgr.Remove(sess.ID)

	if sessMgr.Count() != 0 {
		t.Errorf("Expected 0 sessions after removal, got %d", sessMgr.Count())
	}
}
