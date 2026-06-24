package main

import (
	"fmt"
)

func triggerAiDemo() {
	registry := NewApiProviderRegistry()

	// Register OpenAI
	registry.RegisterProvider("openai", &DummyProvider{Name: "OpenAI-Mock"})
	registry.RegisterModel(&Model{
		ProviderID:     "openai",
		ModelID:        "gpt-4",
		ContextWindow:  8192,
		MaxTokens:      4096,
		SupportsImages: true,
		SupportsTools:  true,
	})

	// Register Anthropic
	registry.RegisterProvider("anthropic", &DummyProvider{Name: "Anthropic-Mock"})
	registry.RegisterModel(&Model{
		ProviderID:     "anthropic",
		ModelID:        "claude-3-opus",
		ContextWindow:  200000,
		MaxTokens:      4096,
		SupportsImages: true,
		SupportsTools:  true,
	})

	fmt.Println("=== Executing AI Stream (OpenAI) ===")
	ctx := &Context{
		SystemPrompt: "You are a helpful coding assistant.",
		Messages: []Message{
			{Role: RoleUser, Content: []ContentPart{{Type: ContentText, Text: "Write a function."}}},
		},
		Tools: []string{"mcp_shell"},
	}

	stream1, _ := registry.ExecuteStream("openai", "gpt-4", ctx)
	for event := range stream1 {
		fmt.Printf("Event: %s - Data: %s\n", event.Type, event.Data)
	}

	fmt.Println("=== Executing AI Stream (Anthropic) ===")
	stream2, _ := registry.ExecuteStream("anthropic", "claude-3-opus", ctx)
	for event := range stream2 {
		fmt.Printf("Event: %s - Data: %s\n", event.Type, event.Data)
	}
}
