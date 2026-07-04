package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// State represents the serialized agent orchestration state.
type State struct {
	SessionID   string    `json:"session_id"`
	LastConnect time.Time `json:"last_connect"`
	Status      string    `json:"status"`
}

// Orchestrator manages the agent's background lifecycle and state synchronization.
type Orchestrator struct {
	cacheDir string
	state    State
	mu       sync.Mutex
}

// NewOrchestrator initializes a new agent orchestrator, ensuring the cache directory exists.
func NewOrchestrator(workspaceRoot string) (*Orchestrator, error) {
	cacheDir := filepath.Join(workspaceRoot, ".pi")

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create agent cache dir: %w", err)
	}

	o := &Orchestrator{
		cacheDir: cacheDir,
		state: State{
			Status: "initialized",
		},
	}

	// Try loading existing state
	_ = o.LoadState()

	return o, nil
}

// SaveState safely persists the current state to the `.pi/state.json` file.
func (o *Orchestrator) SaveState() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	statePath := filepath.Join(o.cacheDir, "state.json")
	data, err := json.MarshalIndent(o.state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal agent state: %w", err)
	}

	if err := os.WriteFile(statePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write agent state: %w", err)
	}

	return nil
}

// LoadState reads the state from the `.pi/state.json` cache file.
func (o *Orchestrator) LoadState() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	statePath := filepath.Join(o.cacheDir, "state.json")
	data, err := os.ReadFile(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No existing state is fine
		}
		return fmt.Errorf("failed to read agent state: %w", err)
	}

	if err := json.Unmarshal(data, &o.state); err != nil {
		return fmt.Errorf("failed to unmarshal agent state: %w", err)
	}

	return nil
}

// ConnectRemote simulates connecting to the remote remote agent orchestration server.
func (o *Orchestrator) ConnectRemote() error {
	o.mu.Lock()
	o.state.Status = "connecting"
	o.mu.Unlock()
	o.SaveState()

	// Simulate connection latency
	time.Sleep(500 * time.Millisecond)

	o.mu.Lock()
	o.state.SessionID = fmt.Sprintf("session-%d", time.Now().Unix())
	o.state.LastConnect = time.Now()
	o.state.Status = "connected"
	o.mu.Unlock()

	return o.SaveState()
}

// Status returns the current connection status in a thread-safe manner.
func (o *Orchestrator) Status() string {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.state.Status
}
