//go:build !windows

package terminal

import (
	"os"
	"os/exec"

	"github.com/creack/pty"
)

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
	s := &LocalSession{pty: f}
	s.SetOpen(true)
	return s, nil
}

func (s *LocalSession) Read(p []byte) (n int, err error)  { return s.pty.Read(p) }
func (s *LocalSession) Write(p []byte) (n int, err error) { return s.pty.Write(p) }
func (s *LocalSession) Close() error {
	s.SetOpen(false)
	return s.pty.Close()
}
func (s *LocalSession) Resize(cols, rows int) error {
	return pty.Setsize(s.pty, &pty.Winsize{Cols: uint16(cols), Rows: uint16(rows)})
}

func (s *LocalSession) GetWorkingDirectory() (string, error) {
	return os.Getwd()
}
