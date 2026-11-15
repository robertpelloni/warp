package workspace

import (
	"sync"
)

// Workspace represents a managed agent session workspace.
type Workspace struct {
	ID           string
	Capabilities []string
}

// Manager orchestrates multiple workspaces (Jules-inspired).
type Manager struct {
	workspaces map[string]*Workspace
	mu         sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		workspaces: make(map[string]*Workspace),
	}
}

func (m *Manager) CreateWorkspace(id string, caps []string) *Workspace {
	m.mu.Lock()
	defer m.mu.Unlock()

	w := &Workspace{
		ID:           id,
		Capabilities: caps,
	}
	m.workspaces[id] = w
	return w
}

func (m *Manager) GetWorkspace(id string) *Workspace {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.workspaces[id]
}
