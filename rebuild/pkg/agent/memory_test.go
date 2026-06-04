package agent

import "testing"

type MockMemoryProvider struct {
	Prompt string
}

func (m *MockMemoryProvider) BuildSystemPrompt() string { return m.Prompt }
func (m *MockMemoryProvider) Sync(u, a string) error    { return nil }

func TestMemoryManager(t *testing.T) {
	mm := NewMemoryManager()
	p := &MockMemoryProvider{Prompt: "Recall previous tasks."}
	mm.AddProvider(p)

	prompt := mm.BuildSystemPrompt()
	if prompt != "Recall previous tasks.\n" {
		t.Errorf("Unexpected prompt: %q", prompt)
	}

	err := mm.SyncAll("hello", "hi")
	if err != nil {
		t.Errorf("SyncAll failed: %v", err)
	}
}
