package agent

import (
	"context"
	"testing"
	"time"

	"github.com/robertpelloni/warp-rebuild/pkg/harness"
	"github.com/robertpelloni/warp-rebuild/pkg/terminal"
)

// Define local mock for integration testing to avoid cross-package test dependency
type TestSession struct {
	terminal.BaseSession
}

func (m *TestSession) Read(p []byte) (n int, err error)  { return 0, nil }
func (m *TestSession) Write(p []byte) (n int, err error) { return 0, nil }
func (m *TestSession) Close() error                     { return nil }
func (m *TestSession) Resize(cols, rows int) error      { return nil }
func (m *TestSession) GetWorkingDirectory() (string, error) { return "", nil }

func TestHarnessIntegration(t *testing.T) {
	// 1. Setup Terminal components
	tm := terminal.NewTabManager()
	bm := terminal.NewBlockManager()

	// Create a mock session for the test
	mockSession := &TestSession{}
	mockSession.SetOpen(true)

	tab := tm.AddTab("Main Tab", mockSession)
	block := bm.CreateBlock("block1", tab.ID, "terminal")

	// 2. Setup Harness components
	registry := harness.NewRegistry()
	model := &harness.ModelInfo{
		ID:            "gpt-4o",
		Provider:      harness.ProviderOpenAI,
		ContextWindow: 128000,
		SupportsTools: true,
	}
	registry.RegisterModel(model)

	// 3. Setup Agent components
	mockProvider := &harness.MockProvider{}
	memory := NewMemoryManager()
	memory.AddProvider(&MockMemoryProvider{Prompt: "Integration Context"})

	cb := NewCircuitBreaker(3, 1*time.Second)

	agent := NewAgent("agent-1", "Standard System Prompt", model, mockProvider)

	// 4. Verification: Ensure all parts are wired and accessible
	if tm.GetActiveTab().ID != tab.ID {
		t.Error("TabManager active tab mismatch")
	}

	if bm.GetBlock("block1").ID != block.ID {
		t.Error("BlockManager block retrieval failure")
	}

	if registry.GetModel("gpt-4o").ID != model.ID {
		t.Error("Harness Registry model retrieval failure")
	}

	// 5. Simulate a simple loop execution through the circuit breaker
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := cb.Execute(func() error {
		return agent.RunLoop(ctx)
	})

	if err != nil {
		t.Errorf("Agent RunLoop via CircuitBreaker failed: %v", err)
	}

	if cb.GetState() != CircuitClosed {
		t.Error("CircuitBreaker state should remain closed on success")
	}

	t.Log("Harness and Go engine integration verified successfully.")
}
