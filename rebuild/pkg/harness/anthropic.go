package harness

import (
	"context"
	"github.com/liushuangls/go-anthropic/v2"
)

type AnthropicProvider struct {
	client *anthropic.Client
	model  string
}

func NewAnthropicProvider(apiKey, model string) *AnthropicProvider {
	return &AnthropicProvider{
		client: anthropic.NewClient(apiKey),
		model:  model,
	}
}

func (p *AnthropicProvider) Chat(ctx context.Context, messages []Message, tools []map[string]interface{}) (*LLMResponse, error) {
	var anthropicMessages []anthropic.Message
	var systemPrompt string

	for _, msg := range messages {
		if msg.Role == "system" {
			systemPrompt = msg.Content
			continue
		}
		anthropicMessages = append(anthropicMessages, anthropic.Message{
			Role:    anthropic.ChatRole(msg.Role),
			Content: []anthropic.MessageContent{anthropic.NewTextMessageContent(msg.Content)},
		})
	}

	resp, err := p.client.CreateMessages(ctx, anthropic.MessagesRequest{
		Model:     anthropic.Model(p.model),
		Messages:  anthropicMessages,
		System:    systemPrompt,
		MaxTokens: 4096,
	})

	if err != nil {
		return nil, err
	}

	res := &LLMResponse{
		Content: *resp.Content[0].Text,
	}

	return res, nil
}
