package terminal

import (
	"io"
	"sync"
)

type Session interface {
	io.ReadWriteCloser
	Resize(cols, rows int) error
	GetWorkingDirectory() (string, error)
}

type BaseSession struct {
	open bool
	mu   sync.RWMutex
}

func (s *BaseSession) IsOpen() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.open
}

func (s *BaseSession) SetOpen(open bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.open = open
}
