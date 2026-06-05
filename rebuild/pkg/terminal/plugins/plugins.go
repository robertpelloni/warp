package plugins

import (
	"fmt"
	"sync"
)

// Plugin defines the interface for terminal extensions.
type Plugin interface {
	Name() string
	Init() error
}

// PluginManager handles loading and life cycle of plugins.
type PluginManager struct {
	plugins map[string]Plugin
	mu      sync.RWMutex
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]Plugin),
	}
}

func (m *PluginManager) LoadPlugin(p Plugin) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := p.Init(); err != nil {
		return fmt.Errorf("plugin %s failed to init: %w", p.Name(), err)
	}

	m.plugins[p.Name()] = p
	return nil
}

func (m *PluginManager) GetPlugins() []Plugin {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list := make([]Plugin, 0, len(m.plugins))
	for _, p := range m.plugins {
		list = append(list, p)
	}
	return list
}
