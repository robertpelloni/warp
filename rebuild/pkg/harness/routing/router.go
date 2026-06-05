package routing

import (
	"context"
	"fmt"
	"github.com/robertpelloni/warp-rebuild/pkg/harness"
)

// Router handles model selection and fallback (OmniRoute-inspired).
type Router struct {
	Primary   harness.Provider
	Fallbacks []harness.Provider
}

func (r *Router) Chat(ctx context.Context, messages []harness.Message, tools []map[string]interface{}) (*harness.LLMResponse, error) {
	resp, err := r.Primary.Chat(ctx, messages, tools)
	if err == nil {
		return resp, nil
	}

	fmt.Printf("Primary failed, trying fallbacks: %v\n", err)
	for _, f := range r.Fallbacks {
		resp, err = f.Chat(ctx, messages, tools)
		if err == nil {
			return resp, nil
		}
	}

	return nil, fmt.Errorf("all providers failed")
}
