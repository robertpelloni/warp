package simulation

import (
	"context"
	"fmt"
	"github.com/robertpelloni/warp-rebuild/pkg/agent"
	"github.com/robertpelloni/warp-rebuild/pkg/harness"
)

type GameTool struct {
	Engine *Engine
}

func (t *GameTool) Name() string        { return "game_action" }
func (t *GameTool) Description() string { return "Performs an action in the game simulation." }
func (t *GameTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"action": map[string]interface{}{
				"type":        "string",
				"description": "The action to take (e.g., move_left, move_right, jump).",
			},
		},
		"required": []string{"action"},
	}
}
func (t *GameTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	action, _ := args["action"].(string)
	t.Engine.Step(0.1)
	return fmt.Sprintf("Action %s executed, simulation stepped.", action), nil
}

func (e *Engine) RunBenchmark(ctx context.Context, iterations int) error {
	model := &harness.ModelInfo{ID: "benchmark-model"}
	mockProvider := &harness.MockProvider{Responses: []*harness.LLMResponse{}}
	for i := 0; i < iterations; i++ {
		mockProvider.Responses = append(mockProvider.Responses, &harness.LLMResponse{
			Content:  "Taking action...",
			ToolCall: &harness.ToolCall{Name: "game_action", Args: map[string]interface{}{"action": "step"}},
		})
	}
	mockProvider.Responses = append(mockProvider.Responses, &harness.LLMResponse{Content: "Benchmark complete."})

	a := agent.NewAgent("benchmark-agent", "Play game", model, mockProvider)
	a.AddTool(&GameTool{Engine: e})

	return a.RunLoop(ctx)
}
