package harness

import "testing"

func TestRegistry(t *testing.T) {
	r := NewRegistry()
	m := &ModelInfo{
		ID:            "gpt-4",
		Provider:      ProviderOpenAI,
		ContextWindow: 128000,
		SupportsTools: true,
	}
	r.RegisterModel(m)

	retrieved := r.GetModel("gpt-4")
	if retrieved == nil {
		t.Fatal("Expected to retrieve model gpt-4")
	}
	if retrieved.ContextWindow != 128000 {
		t.Errorf("Expected context window 128000, got %d", retrieved.ContextWindow)
	}
}
