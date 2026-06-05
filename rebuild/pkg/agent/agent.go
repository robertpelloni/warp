package agent

import (
	"context"
	"github.com/robertpelloni/warp-rebuild/pkg/harness"
)

type Agent struct {
	ID      string
	History []harness.Message
	Provider harness.Provider
	Model   *harness.ModelInfo
}

func NewAgent(id, prompt string, model *harness.ModelInfo, provider harness.Provider) *Agent {
	return &Agent{ID: id, History: []harness.Message{{Role: "system", Content: prompt}}, Provider: provider, Model: model}
}

func (a *Agent) RunLoop(ctx context.Context) error {
	resp, err := a.Provider.Chat(ctx, a.History, nil)
	if err != nil { return err }
	a.History = append(a.History, harness.Message{Role: "assistant", Content: resp.Content})
	return nil
}

type CircuitBreaker struct {
	state int
}
func NewCircuitBreaker(t int, d any) *CircuitBreaker { return &CircuitBreaker{} }
func (cb *CircuitBreaker) Execute(f func() error) error { return f() }
func (cb *CircuitBreaker) GetState() int { return 0 }
const CircuitClosed = 0
