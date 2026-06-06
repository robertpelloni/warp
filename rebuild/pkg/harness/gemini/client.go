package gemini

import (
	"context"
	"fmt"
	"github.com/robertpelloni/warp-rebuild/pkg/harness"
)

// GeminiProvider implements the harness.Provider for Google Gemini.
type GeminiProvider struct {
	APIKey string
	Model  string
}

func NewGeminiProvider(apiKey, model string) *GeminiProvider {
	return &GeminiProvider{APIKey: apiKey, Model: model}
}

func (p *GeminiProvider) Chat(ctx context.Context, messages []harness.Message, tools []map[string]interface{}) (*harness.LLMResponse, error) {
	// For now, we use a structural mock that simulates the Gemini response format.
	// In a real implementation, this would use the "google.golang.org/genai" package.
	if p.APIKey == "" {
		return nil, fmt.Errorf("gemini api key is required")
	}

	fmt.Printf("Gemini Chat calling model %s\n", p.Model)

	// Simulate response for architectural parity
	return &harness.LLMResponse{
		Content: "Gemini: " + messages[len(messages)-1].Content,
	}, nil
}
