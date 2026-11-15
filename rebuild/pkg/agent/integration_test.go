package agent

import (
	"context"
	"testing"
	"time"

	"github.com/robertpelloni/warp-rebuild/pkg/harness"
	"github.com/robertpelloni/warp-rebuild/pkg/terminal"
)

func TestHarnessIntegration(t *testing.T) {
	tm := terminal.NewTabManager()
	bm := terminal.NewBlockManager()
	mockSession := &MockSession{}
	mockSession.SetOpen(true)
	tab := tm.AddTab("Main Tab", mockSession)
	 _ = bm.CreateBlock("block1", tab.ID, "terminal")
	model := &harness.ModelInfo{ID: "gpt-4o"}
	mockProvider := &harness.MockProvider{Responses: []*harness.LLMResponse{{Content: "routed"}}}
	agent := NewAgent("agent-1", "Standard System Prompt", model, mockProvider)
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	_ = agent.RunLoop(ctx)
	if len(agent.History) < 2 { t.Error("History should have assistant response") }
}

type MockSession struct {
	terminal.BaseSession
}
func (m *MockSession) Read(p []byte) (n int, err error)  { return 0, nil }
func (m *MockSession) Write(p []byte) (n int, err error) { return 0, nil }
func (m *MockSession) Close() error                     { return nil }
func (m *MockSession) Resize(cols, rows int) error      { return nil }
func (m *MockSession) GetWorkingDirectory() (string, error) { return "", nil }
