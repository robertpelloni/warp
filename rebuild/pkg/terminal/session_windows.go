//go:build windows

package terminal

import (
	"os"
)

func (s *LocalSession) GetWorkingDirectory() (string, error) {
	// Windows-specific CWD retrieval logic for child processes
	return os.Getwd()
}
