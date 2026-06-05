package agent

import (
	"testing"
	"github.com/robertpelloni/warp-rebuild/pkg/harness"
)

func TestAgent(t *testing.T) {
	m := &harness.ModelInfo{ID: "test"}
	p := &harness.MockProvider{Responses: []*harness.LLMResponse{{Content: "hi"}}}
	a := NewAgent("a1", "sys", m, p)
	if a.ID != "a1" { t.Error("id mismatch") }
}
