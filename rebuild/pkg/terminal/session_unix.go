//go:build !windows

package terminal

import (
	"os"
)

func (s *LocalSession) GetWorkingDirectory() (string, error) {
	// Note: Functional implementation for Unix/Linux would involve
	// reading /proc/[pid]/cwd or using specific platform APIs.
	// This is a placeholder for the future deep PTY integration.
	return os.Getwd()
}
