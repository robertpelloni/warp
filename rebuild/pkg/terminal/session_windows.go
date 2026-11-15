//go:build windows

package terminal

import (
	"errors"
	"os"
)

type LocalSession struct {
	BaseSession
}

func NewLocalSession(shell string) (*LocalSession, error) {
	return nil, errors.New("local pty not yet supported on windows")
}

func (s *LocalSession) Read(p []byte) (n int, err error)  { return 0, os.ErrInvalid }
func (s *LocalSession) Write(p []byte) (n int, err error) { return 0, os.ErrInvalid }
func (s *LocalSession) Close() error                     { return nil }
func (s *LocalSession) Resize(cols, rows int) error      { return nil }

func (s *LocalSession) GetWorkingDirectory() (string, error) {
	return os.Getwd()
}
