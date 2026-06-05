package e2e

import (
	"context"
	"testing"
	"time"
	"github.com/robertpelloni/warp-rebuild/pkg/agent"
	"github.com/robertpelloni/warp-rebuild/pkg/agent/service"
	"github.com/robertpelloni/warp-rebuild/pkg/harness"
)

func TestEndToEndAgentSession(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	svc := service.NewAgentService(10001)
	go svc.Start(ctx)
	time.Sleep(100 * time.Millisecond)
	model := &harness.ModelInfo{ID: "e2e-model"}
	mockProvider := &harness.MockProvider{Responses: []*harness.LLMResponse{{Content: "E2E Resp"}}}
	a := agent.NewAgent("e2e-agent", "prompt", model, mockProvider)
	_ = a.RunLoop(ctx)
	if !svc.IsActive() { t.Error("Service should be active") }
	_ = svc.Stop(ctx)
}
