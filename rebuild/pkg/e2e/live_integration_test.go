package e2e

import (
	"context"
	"os"
	"testing"
	"github.com/robertpelloni/warp-rebuild/pkg/harness"
	"github.com/robertpelloni/warp-rebuild/pkg/agent"
)

func TestLiveLLMIntegration(t *testing.T) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping live integration test: ANTHROPIC_API_KEY not set")
	}

	provider := harness.NewAnthropicProvider(apiKey, "claude-3-5-sonnet-20240620")
	a := agent.NewAgent("test", "You are a concise assistant.", &harness.ModelInfo{ID: "claude-3-5-sonnet"}, provider)
	a.History = append(a.History, harness.Message{Role: "user", Content: "Say 'Hello, World!'"})

	err := a.RunLoop(context.Background())
	if err != nil {
		t.Fatalf("Failed to run agent loop: %v", err)
	}

	lastMsg := a.History[len(a.History)-1]
	t.Logf("Agent response: %s", lastMsg.Content)
	if lastMsg.Role != "assistant" {
		t.Errorf("Expected role 'assistant', got '%s'", lastMsg.Role)
	}
}
