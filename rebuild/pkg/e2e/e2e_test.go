package e2e

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/robertpelloni/warp-rebuild/pkg/agent"
	"github.com/robertpelloni/warp-rebuild/pkg/agent/service"
	"github.com/robertpelloni/warp-rebuild/pkg/harness"
	"github.com/robertpelloni/warp-rebuild/pkg/terminal"
)

func TestEndToEndAgentSession(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Create a dummy file for the agent to read
	dummyFile, _ := os.CreateTemp("", "e2e-test-file")
	defer os.Remove(dummyFile.Name())
	dummyFile.WriteString("E2E Secret Content")
	dummyFile.Close()

	// 2. Setup Mock Provider with a tool-calling sequence
	mockProvider := &harness.MockProvider{
		Responses: []*harness.LLMResponse{
			{
				Content: "I will read the test file.",
				ToolCall: &harness.ToolCall{
					Name: "read_file",
					Args: map[string]interface{}{"path": dummyFile.Name()},
				},
			},
			{
				Content: "I have read the file. The content is E2E Secret Content. Now I will finish.",
			},
		},
	}

	// 3. Start Agent Service
	svc := service.NewAgentService(10001)
	go func() {
		_ = svc.Start(ctx)
	}()
	time.Sleep(200 * time.Millisecond)

	// 4. Initialize Terminal Session
	session, err := terminal.NewLocalSession("echo")
	if err != nil {
		t.Fatalf("Failed to start terminal session: %v", err)
	}
	defer session.Close()

	// 5. Setup Agent and Harness
	model := &harness.ModelInfo{ID: "e2e-model"}
	a := agent.NewAgent("e2e-agent", "You are an E2E test agent.", model, mockProvider)

	// 6. Run Agentic Loop
	err = a.RunLoop(ctx)
	if err != nil {
		t.Errorf("Agent loop failed in E2E session: %v", err)
	}

	// 7. Verify the agent's history contains the observation
	foundObservation := false
	for _, msg := range a.History {
		if msg.Role == "user" && msg.Content == "Observation: E2E Secret Content" {
			foundObservation = true
			break
		}
	}
	if !foundObservation {
		t.Error("Agent history did not contain the expected file observation")
	}

	// 8. Verify Service Status
	if !svc.IsActive() {
		t.Error("Agent service should be active during session")
	}

	// 9. Cleanup
	_ = svc.Stop(ctx)

	t.Log("End-to-End agentic session with tool-use verified successfully.")
}
