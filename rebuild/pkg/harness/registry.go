package harness

import (
	"time"
)

type ProviderType string

const (
	ProviderAnthropic  ProviderType = "anthropic"
	ProviderOpenAI     ProviderType = "openai"
	ProviderGoogle     ProviderType = "google"
	ProviderOpenRouter ProviderType = "openrouter"
)

type ProviderConfig struct {
	Type    ProviderType  `json:"type"`
	APIKey  string        `json:"apiKey"`
	BaseURL string        `json:"baseURL,omitempty"`
	Timeout time.Duration `json:"timeout"`
}

type ModelInfo struct {
	ID            string       `json:"id"`
	Provider      ProviderType `json:"provider"`
	ContextWindow int          `json:"contextWindow"`
	SupportsTools bool         `json:"supportsTools"`
}

type Registry struct {
	providers map[ProviderType]*ProviderConfig
	models    map[string]*ModelInfo
}

func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[ProviderType]*ProviderConfig),
		models:    make(map[string]*ModelInfo),
	}
}

func (r *Registry) RegisterProvider(config *ProviderConfig) {
	r.providers[config.Type] = config
}

func (r *Registry) RegisterModel(model *ModelInfo) {
	r.models[model.ID] = model
}

func (r *Registry) GetModel(id string) *ModelInfo {
	return r.models[id]
}
