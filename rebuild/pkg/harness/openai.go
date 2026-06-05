package harness

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

type OpenAIProvider struct {
	client *openai.Client
	model  string
}

func NewOpenAIProvider(apiKey, model string) *OpenAIProvider {
	return &OpenAIProvider{client: openai.NewClient(apiKey), model: model}
}

func (p *OpenAIProvider) Chat(ctx context.Context, messages []Message, tools []map[string]interface{}) (*LLMResponse, error) {
	// Minimal shim for restoration
	return &LLMResponse{Content: "OpenAI Resp"}, nil
}

type MockProvider struct {
	Responses []*LLMResponse
	idx       int
}

func (p *MockProvider) Chat(ctx context.Context, messages []Message, tools []map[string]interface{}) (*LLMResponse, error) {
	if p.idx >= len(p.Responses) { return &LLMResponse{Content: "End."}, nil }
	res := p.Responses[p.idx]
	p.idx++
	return res, nil
}
