package harness

import (
	"context"
	"encoding/json"

	"github.com/sashabaranov/go-openai"
)

type OpenAIProvider struct {
	client *openai.Client
	model  string
}

func NewOpenAIProvider(apiKey, model string) *OpenAIProvider {
	return &OpenAIProvider{
		client: openai.NewClient(apiKey),
		model:  model,
	}
}

func (p *OpenAIProvider) Chat(ctx context.Context, messages []Message, tools []map[string]interface{}) (*LLMResponse, error) {
	apiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, m := range messages {
		apiMessages[i] = openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}

	apiTools := make([]openai.Tool, len(tools))
	for i, t := range tools {
		apiTools[i] = openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:       t["name"].(string),
				Description: t["description"].(string),
				Parameters: t["parameters"],
			},
		}
	}

	resp, err := p.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    p.model,
		Messages: apiMessages,
		Tools:    apiTools,
	})
	if err != nil {
		return nil, err
	}

	result := &LLMResponse{
		Content: resp.Choices[0].Message.Content,
	}

	if len(resp.Choices[0].Message.ToolCalls) > 0 {
		tc := resp.Choices[0].Message.ToolCalls[0]
		var args map[string]interface{}
		_ = json.Unmarshal([]byte(tc.Function.Arguments), &args)
		result.ToolCall = &ToolCall{
			Name: tc.Function.Name,
			Args: args,
		}
	}

	return result, nil
}

// MockProvider for testing
type MockProvider struct {
	Responses []*LLMResponse
	idx       int
}

func (p *MockProvider) Chat(ctx context.Context, messages []Message, tools []map[string]interface{}) (*LLMResponse, error) {
	if p.idx >= len(p.Responses) {
		return &LLMResponse{Content: "End of conversation."}, nil
	}
	res := p.Responses[p.idx]
	p.idx++
	return res, nil
}
