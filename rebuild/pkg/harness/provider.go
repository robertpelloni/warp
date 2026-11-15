package harness

import "context"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LLMResponse struct {
	Content string
	ToolCall *ToolCall
}

type ToolCall struct {
	Name string
	Args map[string]interface{}
}

type Provider interface {
	Chat(ctx context.Context, messages []Message, tools []map[string]interface{}) (*LLMResponse, error)
}

type ModelInfo struct {
	ID string
}
