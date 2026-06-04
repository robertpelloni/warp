package terminal

import (
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/creack/pty"
)

// Session defines the interface for a terminal session.
type Session interface {
	io.ReadWriteCloser
	Resize(cols, rows int) error
	GetWorkingDirectory() (string, error)
}

// BaseSession provides a common foundation for terminal sessions.
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

// LocalSession implements Session for a local PTY.
type LocalSession struct {
	BaseSession
	pty *os.File
}

func NewLocalSession(shell string) (*LocalSession, error) {
	c := exec.Command(shell)
	f, err := pty.Start(c)
	if err != nil {
		return nil, err
	}

	s := &LocalSession{
		pty: f,
	}
	s.SetOpen(true)
	return s, nil
}

func (s *LocalSession) Read(p []byte) (n int, err error) {
	return s.pty.Read(p)
}

func (s *LocalSession) Write(p []byte) (n int, err error) {
	return s.pty.Write(p)
}

func (s *LocalSession) Close() error {
	s.SetOpen(false)
	return s.pty.Close()
}

func (s *LocalSession) Resize(cols, rows int) error {
	return pty.Setsize(s.pty, &pty.Winsize{
		Cols: uint16(cols),
		Rows: uint16(rows),
	})
}

func (s *LocalSession) GetWorkingDirectory() (string, error) {
	// Simple implementation: this usually requires OS-specific logic (e.g. reading from /proc)
	return os.Getwd()
}
