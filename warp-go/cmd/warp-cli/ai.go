package main

import (
	"fmt"
	"strings"
	"sync"
)

// --- Types ---

// MessageRole represents the role of a message in a conversation.
type MessageRole string

const (
	RoleSystem    MessageRole = "system"
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleTool      MessageRole = "tool"
)

// Message content components
type ContentPartType string

const (
	ContentText     ContentPartType = "text"
	ContentImage    ContentPartType = "image"
	ContentToolCall ContentPartType = "tool_call"
)

type ContentPart struct {
	Type   ContentPartType
	Text   string // If type is text
	CallID string // If type is tool_call
	Name   string // If type is tool_call
	Args   string // If type is tool_call JSON args
}

// Message represents a single turn in a conversation context.
type Message struct {
	Role    MessageRole
	Content []ContentPart
}

// Context represents the context payload passed to a provider.
type Context struct {
	SystemPrompt string
	Messages     []Message
	Tools        []string // Just string descriptions/names for now
}

// Model represents a specific model string within an API provider.
type Model struct {
	ProviderID     string
	ModelID        string
	ContextWindow  int
	MaxTokens      int
	CostPer1kIn    float64
	CostPer1kOut   float64
	SupportsImages bool
	SupportsTools  bool
}

// --- Stream Events ---

type EventType string

const (
	EventStart         EventType = "start"
	EventTextDelta     EventType = "text_delta"
	EventToolCallStart EventType = "tool_call_start"
	EventToolCallDelta EventType = "tool_call_delta"
	EventDone          EventType = "done"
	EventError         EventType = "error"
)

type AssistantMessageEvent struct {
	Type EventType
	Data string
	Err  error
}

type AssistantMessageEventStream <-chan AssistantMessageEvent

// --- Api Provider interface ---

type ApiProvider interface {
	Stream(model *Model, ctx *Context) AssistantMessageEventStream
}

// --- Registry ---

type ApiProviderRegistry struct {
	providers map[string]ApiProvider
	models    map[string]*Model
	mu        sync.RWMutex
}

func NewApiProviderRegistry() *ApiProviderRegistry {
	return &ApiProviderRegistry{
		providers: make(map[string]ApiProvider),
		models:    make(map[string]*Model),
	}
}

func (r *ApiProviderRegistry) RegisterProvider(id string, provider ApiProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[id] = provider
	fmt.Printf("[ApiRegistry] Registered provider: %s\n", id)
}

func (r *ApiProviderRegistry) RegisterModel(model *Model) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := fmt.Sprintf("%s:%s", model.ProviderID, model.ModelID)
	r.models[key] = model
	fmt.Printf("[ApiRegistry] Registered model: %s\n", key)
}

func (r *ApiProviderRegistry) GetModel(providerID, modelID string) (*Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := fmt.Sprintf("%s:%s", providerID, modelID)
	model, exists := r.models[key]
	if !exists {
		return nil, fmt.Errorf("model %s not found", key)
	}
	return model, nil
}

func (r *ApiProviderRegistry) GetProvider(providerID string) (ApiProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, exists := r.providers[providerID]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", providerID)
	}
	return p, nil
}

func (r *ApiProviderRegistry) ExecuteStream(providerID, modelID string, ctx *Context) (AssistantMessageEventStream, error) {
	model, err := r.GetModel(providerID, modelID)
	if err != nil {
		return nil, err
	}
	provider, err := r.GetProvider(providerID)
	if err != nil {
		return nil, err
	}

	return provider.Stream(model, ctx), nil
}

// --- Dummy Provider Impl ---

type DummyProvider struct {
	Name string
}

func (p *DummyProvider) Stream(model *Model, ctx *Context) AssistantMessageEventStream {
	stream := make(chan AssistantMessageEvent)

	go func() {
		defer close(stream)
		stream <- AssistantMessageEvent{Type: EventStart}

		// Simulate reasoning/text output based on context length
		words := strings.Split("This is a simulated streaming response from "+p.Name+" using model "+model.ModelID+". ", " ")
		for _, w := range words {
			stream <- AssistantMessageEvent{Type: EventTextDelta, Data: w + " "}
		}

		// If tools were provided, simulate a tool call
		if len(ctx.Tools) > 0 {
			stream <- AssistantMessageEvent{Type: EventToolCallStart, Data: "mcp_shell"}
			stream <- AssistantMessageEvent{Type: EventToolCallDelta, Data: `{"cmd": "echo hello"}`}
		}

		stream <- AssistantMessageEvent{Type: EventDone}
	}()

	return stream
}
