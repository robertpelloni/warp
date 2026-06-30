//go:build !windows

package pty

import "fmt"

func (p *PTY) startWindows(cfg Config) error {
	return fmt.Errorf("not implemented on this platform")
}
