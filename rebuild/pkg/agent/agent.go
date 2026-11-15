package agent

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/robertpelloni/warp-rebuild/pkg/harness"
)

var ErrCircuitOpen = errors.New("circuit breaker is open")

type Agent struct {
	ID       string
	History  []harness.Message
	Provider harness.Provider
	Model    *harness.ModelInfo
	Tools    map[string]Tool
}

func NewAgent(id, prompt string, model *harness.ModelInfo, provider harness.Provider) *Agent {
	return &Agent{
		ID:       id,
		History:  []harness.Message{{Role: "system", Content: prompt}},
		Provider: provider,
		Model:    model,
		Tools:    make(map[string]Tool),
	}
}

func (a *Agent) AddTool(t Tool) {
	a.Tools[t.Name()] = t
}

func (a *Agent) RunLoop(ctx context.Context) error {
	var toolDefs []map[string]interface{}
	for _, t := range a.Tools {
		toolDefs = append(toolDefs, map[string]interface{}{
			"name":        t.Name(),
			"description": t.Description(),
			"parameters":  t.Parameters(),
		})
	}

	for {
		resp, err := a.Provider.Chat(ctx, a.History, toolDefs)
		if err != nil {
			return err
		}

		a.History = append(a.History, harness.Message{Role: "assistant", Content: resp.Content})

		if resp.ToolCall == nil {
			break
		}

		tool, ok := a.Tools[resp.ToolCall.Name]
		var observation string
		if !ok {
			observation = "Error: tool not found"
		} else {
			observation, err = tool.Execute(ctx, resp.ToolCall.Args)
			if err != nil {
				observation = "Error: " + err.Error()
			}
		}

		a.History = append(a.History, harness.Message{Role: "tool", Content: observation})
		// Loop continues to let LLM process the tool output
	}
	return nil
}

const (
	CircuitClosed = iota
	CircuitOpen
	CircuitHalfOpen
)

type CircuitBreaker struct {
	mutex           sync.Mutex
	state           int
	failureCount    int
	threshold       int
	resetTimeout    time.Duration
	lastFailureTime time.Time
}

func NewCircuitBreaker(threshold int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		threshold:    threshold,
		resetTimeout: resetTimeout,
		state:        CircuitClosed,
	}
}

func (cb *CircuitBreaker) GetState() int {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	return cb.currentState()
}

func (cb *CircuitBreaker) currentState() int {
	if cb.state == CircuitOpen && time.Since(cb.lastFailureTime) > cb.resetTimeout {
		return CircuitHalfOpen
	}
	return cb.state
}

func (cb *CircuitBreaker) Execute(f func() error) error {
	cb.mutex.Lock()
	state := cb.currentState()
	if state == CircuitOpen {
		cb.mutex.Unlock()
		return ErrCircuitOpen
	}
	cb.mutex.Unlock()

	err := f()

	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if err != nil {
		cb.failureCount++
		cb.lastFailureTime = time.Now()
		if state == CircuitHalfOpen || cb.failureCount >= cb.threshold {
			cb.state = CircuitOpen
		}
		return err
	}

	if state == CircuitHalfOpen {
		cb.state = CircuitClosed
		cb.failureCount = 0
	}
	return nil
}
