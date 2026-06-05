package routing

import (
	"context"
	"fmt"
	"testing"
	"github.com/robertpelloni/warp-rebuild/pkg/harness"
)

type FailProvider struct{}
func (p *FailProvider) Chat(ctx context.Context, m []harness.Message, t []map[string]interface{}) (*harness.LLMResponse, error) {
	return nil, fmt.Errorf("fail")
}

func TestRouter(t *testing.T) {
	mockOk := &harness.MockProvider{Responses: []*harness.LLMResponse{{Content: "ok"}}}
	r := &Router{
		Primary:   &FailProvider{},
		Fallbacks: []harness.Provider{mockOk},
	}

	res, err := r.Chat(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("Router failed: %v", err)
	}
	if res.Content != "ok" {
		t.Errorf("Expected fallback response 'ok', got %q", res.Content)
	}
}
