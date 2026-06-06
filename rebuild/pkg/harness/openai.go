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

	resp, err := p.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    p.model,
		Messages: oaMessages,
	})

	if err != nil {
		return nil, err
	}

	return &LLMResponse{
		Content: resp.Choices[0].Message.Content,
	}, nil
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
