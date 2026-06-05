package agent

import (
	"context"
	"fmt"

	"github.com/robertpelloni/warp-rebuild/pkg/agent/tools"
	"github.com/robertpelloni/warp-rebuild/pkg/harness"
)

// Agent represents an AI agent.
type Agent struct {
	ID           string
	SystemPrompt string
	Model        *harness.ModelInfo
	Provider     harness.Provider
	Tools        map[string]tools.Tool
	History      []harness.Message
}

func NewAgent(id, systemPrompt string, model *harness.ModelInfo, provider harness.Provider) *Agent {
	a := &Agent{
		ID:           id,
		SystemPrompt: systemPrompt,
		Model:        model,
		Provider:     provider,
		Tools:        make(map[string]tools.Tool),
		History:      []harness.Message{{Role: "system", Content: systemPrompt}},
	}
	// Add default tools
	a.AddTool(&tools.ReadFilesTool{})
	a.AddTool(&tools.ExecuteCommandTool{})
	return a
}

func (a *Agent) AddTool(t tools.Tool) {
	def := t.GetDefinition()
	name := def["name"].(string)
	a.Tools[name] = t
}

// RunLoop starts the agent's autonomous loop.
func (a *Agent) RunLoop(ctx context.Context) error {
	fmt.Printf("Agent %s starting loop with model %s\n", a.ID, a.Model.ID)

	toolDefs := make([]map[string]interface{}, 0, len(a.Tools))
	for _, t := range a.Tools {
		toolDefs = append(toolDefs, t.GetDefinition())
	}

	for i := 0; i < 5; i++ { // Limit to 5 turns for now
		resp, err := a.Provider.Chat(ctx, a.History, toolDefs)
		if err != nil {
			return err
		}

		if resp.Content != "" {
			fmt.Printf("Agent Thinking: %s\n", resp.Content)
			a.History = append(a.History, harness.Message{Role: "assistant", Content: resp.Content})
		}

		if resp.ToolCall != nil {
			fmt.Printf("Agent Action: %s(%v)\n", resp.ToolCall.Name, resp.ToolCall.Args)
			t, ok := a.Tools[resp.ToolCall.Name]
			if !ok {
				return fmt.Errorf("tool not found: %s", resp.ToolCall.Name)
			}
			result, err := t.Execute(resp.ToolCall.Args)
			if err != nil {
				result = fmt.Sprintf("Error: %v", err)
			}
			fmt.Printf("Agent Observation: %s\n", result)
			a.History = append(a.History, harness.Message{Role: "user", Content: fmt.Sprintf("Observation: %s", result)})
		} else {
			// No more tool calls, turn complete
			break
		}
	}

	return nil
}
