package service

import (
	"context"
	"testing"
	"time"
	"github.com/robertpelloni/warp-rebuild/pkg/harness"
)

func TestService(t *testing.T) {
	s := NewAgentService(10001, &harness.MockProvider{Responses: []*harness.LLMResponse{{Content: "Mock"}}})
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	go s.Start(ctx)
	time.Sleep(100 * time.Millisecond)
	if !s.IsActive() { t.Error("not active") }
	s.Stop(ctx)
}
