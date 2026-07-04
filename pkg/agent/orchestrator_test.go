package agent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOrchestrator_Lifecycle(t *testing.T) {
	// Setup a temporary workspace
	tmpDir, err := os.MkdirTemp("", "warp-agent-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test Initialization
	o, err := NewOrchestrator(tmpDir)
	if err != nil {
		t.Fatalf("NewOrchestrator failed: %v", err)
	}

	if o.Status() != "initialized" {
		t.Errorf("Expected status 'initialized', got '%s'", o.Status())
	}

	// Verify .pi directory was created
	piDir := filepath.Join(tmpDir, ".pi")
	if info, err := os.Stat(piDir); err != nil || !info.IsDir() {
		t.Errorf("Expected .pi directory to exist")
	}

	// Test State Saving
	err = o.SaveState()
	if err != nil {
		t.Errorf("SaveState failed: %v", err)
	}

	stateFile := filepath.Join(piDir, "state.json")
	if _, err := os.Stat(stateFile); err != nil {
		t.Errorf("Expected state.json to exist")
	}

	// Test Remote Connection Simulation
	err = o.ConnectRemote()
	if err != nil {
		t.Errorf("ConnectRemote failed: %v", err)
	}

	if o.Status() != "connected" {
		t.Errorf("Expected status 'connected', got '%s'", o.Status())
	}

	if o.state.SessionID == "" {
		t.Errorf("Expected SessionID to be set after connection")
	}

	// Test State Loading
	o2, err := NewOrchestrator(tmpDir)
	if err != nil {
		t.Fatalf("NewOrchestrator failed on reload: %v", err)
	}

	if o2.Status() != "connected" {
		t.Errorf("Expected loaded status 'connected', got '%s'", o2.Status())
	}

	if o2.state.SessionID != o.state.SessionID {
		t.Errorf("Expected loaded SessionID '%s', got '%s'", o.state.SessionID, o2.state.SessionID)
	}
}
