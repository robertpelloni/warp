package gemini

import (
	"context"
	"testing"
)

func TestGeminiProvider(t *testing.T) {
	p := NewGeminiProvider("key", "gemini-1.5-flash")
	res, err := p.Chat(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("Chat failed: %v", err)
	}
	if res.Content != "Gemini implementation ready." {
		t.Errorf("Unexpected response content: %q", res.Content)
	}
}
