package session

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/robertpelloni/warp/pkg/terminal"
)

// Session represents a terminal session (tab).
type Session struct {
	ID       string
	Name     string
	Terminal *terminal.Terminal
	Active   bool
	Icon     string
}

// Manager handles terminal sessions/tabs.
type Manager struct {
	mu       sync.RWMutex
	sessions []*Session
	active   string // ID of active session
}

// NewManager creates a new session manager.
func NewManager() *Manager {
	return &Manager{
		sessions: make([]*Session, 0),
	}
}

// Create makes a new terminal session and starts it.
func (m *Manager) Create(name string, shell string, cols, rows int, onUpdate func()) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := uuid.New().String()[:8]
	if name == "" {
		name = fmt.Sprintf("Terminal %d", len(m.sessions)+1)
	}

	cfg := terminal.Config{
		Shell:    shell,
		Cols:     cols,
		Rows:     rows,
		OnUpdate: onUpdate,
	}
	term := terminal.New(cfg)

	sess := &Session{
		ID:       id,
		Name:     name,
		Terminal: term,
		Active:   true,
		Icon:     ">_",
	}

	// Deactivate other sessions
	for _, s := range m.sessions {
		s.Active = false
	}

	m.sessions = append(m.sessions, sess)
	m.active = id

	return sess, nil
}

// Get returns a session by ID.
func (m *Manager) Get(id string) *Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, s := range m.sessions {
		if s.ID == id {
			return s
		}
	}
	return nil
}

// Active returns the active session.
func (m *Manager) Active() *Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, s := range m.sessions {
		if s.ID == m.active {
			return s
		}
	}
	return nil
}

// SetActive switches to a session.
func (m *Manager) SetActive(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, s := range m.sessions {
		s.Active = s.ID == id
	}
	m.active = id
}

// List returns all sessions.
func (m *Manager) List() []*Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]*Session, len(m.sessions))
	copy(result, m.sessions)
	return result
}

// Remove closes and removes a session.
func (m *Manager) Remove(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, s := range m.sessions {
		if s.ID == id {
			s.Terminal.Stop()
			m.sessions = append(m.sessions[:i], m.sessions[i+1:]...)
			if m.active == id && len(m.sessions) > 0 {
				m.sessions[0].Active = true
				m.active = m.sessions[0].ID
			}
			return
		}
	}
}

// Count returns the number of sessions.
func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.sessions)
}
