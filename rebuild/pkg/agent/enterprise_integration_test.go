package agent

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/robertpelloni/warp-rebuild/pkg/agent/skills"
	"github.com/robertpelloni/warp-rebuild/pkg/agent/supervisors"
	"github.com/robertpelloni/warp-rebuild/pkg/harness"
)

func TestEnterpriseOrchestrationIntegration(t *testing.T) {
	// 1. Setup Harness
	model := &harness.ModelInfo{ID: "enterprise-model"}
	mockProvider := &harness.MockProvider{
		Responses: []*harness.LLMResponse{
			{Content: "Skill check requested."},
			{Content: "Executing skill procedure."},
		},
	}

	// 2. Setup Skills
	sm := skills.NewSkillManager()
	sm.Register(&skills.Skill{
		Name:        "enterprise-audit",
		Description: "Audit system security.",
	})

	// 3. Setup Supervisor with custom routing
	sup, err := supervisors.CreateSupervisor("default", mockProvider)
	if err != nil {
		t.Fatalf("Failed to create supervisor: %v", err)
	}

	// 4. Setup Agent with Skill Awareness
	// For integration, we append skills to the system prompt
	systemPrompt := "Standard prompt. " + sm.BuildPrompt()
	a := NewAgent("enterprise-agent", systemPrompt, model, mockProvider)

	// 5. Verification: Test routing via supervisor
	supRes, err := sup.Route("Verify enterprise skills")
	if err != nil {
		t.Errorf("Supervisor routing failed: %v", err)
	}
	if supRes != "Skill check requested." {
		t.Errorf("Unexpected supervisor response: %q", supRes)
	}

	// 6. Verification: Run agent loop with integrated context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = a.RunLoop(ctx)
	if err != nil {
		t.Errorf("Integrated agent loop failed: %v", err)
	}

	// Final check on history to ensure skill context was present in system prompt
	if !strings.Contains(a.History[0].Content, "enterprise-audit") {
		t.Error("Agent history system prompt missing registered skill context")
	}

	t.Log("Enterprise orchestration integration verified successfully.")
}
