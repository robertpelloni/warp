package terminal

import (
	"fmt"
	"sync"
)

type Tab struct {
	ID      string
	Title   string
	Session Session
}

type TabManager struct {
	tabs   []*Tab
	active int
	mu     sync.RWMutex
}

func NewTabManager() *TabManager {
	return &TabManager{tabs: make([]*Tab, 0), active: -1}
}

func (m *TabManager) AddTab(title string, session Session) *Tab {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := fmt.Sprintf("tab-%d", len(m.tabs))
	tab := &Tab{ID: id, Title: title, Session: session}
	m.tabs = append(m.tabs, tab)
	if m.active == -1 { m.active = 0 }
	return tab
}

func (m *TabManager) GetActiveTab() *Tab {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.active == -1 { return nil }
	return m.tabs[m.active]
}
