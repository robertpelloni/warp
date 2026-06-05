package agent

import (
	"context"
	"testing"
	"github.com/robertpelloni/warp-rebuild/pkg/harness"
)

func TestAgentCreation(t *testing.T) {
	model := &harness.ModelInfo{ID: "test-model"}
	mockProvider := &harness.MockProvider{}
	a := NewAgent("agent1", "You are a helpful assistant.", model, mockProvider)

	if a.ID != "agent1" {
		t.Errorf("Expected ID agent1, got %s", a.ID)
	}

	ctx := context.Background()
	err := a.RunLoop(ctx)
	if err != nil {
		t.Errorf("RunLoop failed: %v", err)
	}
}
