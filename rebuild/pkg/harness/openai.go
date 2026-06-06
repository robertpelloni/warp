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
	var oaMessages []openai.ChatCompletionMessage
	for _, msg := range messages {
		oaMessages = append(oaMessages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	var oaTools []openai.Tool
	for _, t := range tools {
		name, _ := t["name"].(string)
		desc, _ := t["description"].(string)
		params, _ := t["parameters"].(map[string]interface{})
		oaTools = append(oaTools, openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        name,
				Description: desc,
				Parameters:  params,
			},
		})
	}

	resp, err := p.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    p.model,
		Messages: oaMessages,
		Tools:    oaTools,
	})

	if err != nil {
		return nil, err
	}

	res := &LLMResponse{
		Content: resp.Choices[0].Message.Content,
	}

	if len(resp.Choices[0].Message.ToolCalls) > 0 {
		tc := resp.Choices[0].Message.ToolCalls[0]
		res.ToolCall = &ToolCall{
			Name: tc.Function.Name,
			Args: make(map[string]interface{}), // Simplified: would parse JSON in production
		}
	}

	return res, nil
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
