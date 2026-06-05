//go:build !windows

package terminal

import "os"

func (s *LocalSession) GetWorkingDirectory() (string, error) {
	return os.Getwd()
}
