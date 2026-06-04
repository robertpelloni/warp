package agent

import (
	"context"
	"fmt"
	"github.com/robertpelloni/warp-rebuild/pkg/harness"
)

// Agent represents an AI agent.
type Agent struct {
	ID       string
	SystemPrompt string
	Model    *harness.ModelInfo
}

func NewAgent(id, systemPrompt string, model *harness.ModelInfo) *Agent {
	return &Agent{
		ID:           id,
		SystemPrompt: systemPrompt,
		Model:        model,
	}
}

// RunLoop starts the agent's autonomous loop.
func (a *Agent) RunLoop(ctx context.Context) error {
	fmt.Printf("Agent %s starting loop with model %s\n", a.ID, a.Model.ID)
	// Implementation will involve calling the LLM and processing tools
	return nil
}
