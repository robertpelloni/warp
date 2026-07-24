//go:build !windows

package pty

func (p *PTY) startWindows(cfg Config) error {
	return fmt.Errorf("not implemented on this platform")
}
