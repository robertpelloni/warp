package harness

import (
	"context"
	"encoding/json"
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

	var anthropicTools []anthropic.ToolDefinition
	for _, t := range tools {
		name, _ := t["name"].(string)
		desc, _ := t["description"].(string)
		params, _ := t["parameters"].(map[string]interface{})
		anthropicTools = append(anthropicTools, anthropic.ToolDefinition{
			Name:        name,
			Description: desc,
			InputSchema: params,
		})
	}

	req := anthropic.MessagesRequest{
		Model:     anthropic.Model(p.model),
		Messages:  anthropicMessages,
		System:    systemPrompt,
		MaxTokens: 4096,
	}
	if len(anthropicTools) > 0 {
		req.Tools = anthropicTools
	}

	resp, err := p.client.CreateMessages(ctx, req)

	if err != nil {
		return nil, err
	}

	res := &LLMResponse{}
	for _, c := range resp.Content {
		if c.Type == anthropic.MessagesContentTypeText {
			res.Content = *c.Text
		}
		if c.Type == anthropic.MessagesContentTypeToolUse {
			var args map[string]interface{}
			_ = json.Unmarshal(c.MessageContentToolUse.Input, &args)
			res.ToolCall = &ToolCall{
				Name: c.MessageContentToolUse.Name,
				Args: args,
			}
		}
	}

	return res, nil
}
