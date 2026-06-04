package agent

import (
	"sync"
)

// MemoryProvider defines the interface for a memory source.
type MemoryProvider interface {
	BuildSystemPrompt() string
	Sync(userMsg, assistantResponse string) error
}

// MemoryManager orchestrates memory providers.
type MemoryManager struct {
	providers []MemoryProvider
	mu        sync.RWMutex
}

func NewMemoryManager() *MemoryManager {
	return &MemoryManager{
		providers: make([]MemoryProvider, 0),
	}
}

func (m *MemoryManager) AddProvider(p MemoryProvider) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.providers = append(m.providers, p)
}

func (m *MemoryManager) BuildSystemPrompt() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var prompt string
	for _, p := range m.providers {
		prompt += p.BuildSystemPrompt() + "\n"
	}
	return prompt
}

func (m *MemoryManager) SyncAll(userMsg, assistantResponse string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, p := range m.providers {
		if err := p.Sync(userMsg, assistantResponse); err != nil {
			return err
		}
	}
	return nil
}
