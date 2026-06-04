package terminal

import (
	"fmt"
	"sync"
)

// Tab represents a terminal tab.
type Tab struct {
	ID      string
	Title   string
	Session Session
}

// TabManager manages multiple terminal tabs.
type TabManager struct {
	tabs   []*Tab
	active int
	mu     sync.RWMutex
}

func NewTabManager() *TabManager {
	return &TabManager{
		tabs:   make([]*Tab, 0),
		active: -1,
	}
}

func (m *TabManager) AddTab(title string, session Session) *Tab {
	m.mu.Lock()
	defer m.mu.Unlock()

	tab := &Tab{
		ID:      fmt.Sprintf("%d", len(m.tabs)), // Simple ID for now
		Title:   title,
		Session: session,
	}
	m.tabs = append(m.tabs, tab)
	if m.active == -1 {
		m.active = 0
	}
	return tab
}

func (m *TabManager) GetTabs() []*Tab {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tabs
}

func (m *TabManager) GetActiveTab() *Tab {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.active == -1 || m.active >= len(m.tabs) {
		return nil
	}
	return m.tabs[m.active]
}

func (m *TabManager) SetActiveTab(index int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if index >= 0 && index < len(m.tabs) {
		m.active = index
	}
}
