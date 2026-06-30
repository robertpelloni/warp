package pty

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"
)

// PTY wraps a pseudo-terminal with cross-platform support.
// On Windows it uses ConPTY; on Unix it uses a piped subprocess.
type PTY struct {
	mu      sync.Mutex
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.Reader
	running bool
	wsCol   uint16
	wsRow   uint16

	// Platform-specific
	platform interface{}
}

type platformPTY interface {
	Write(data []byte) (int, error)
	Read(buf []byte) (int, error)
	Resize(cols, rows uint16) error
	Close() error
}

// Config for creating a PTY.
type Config struct {
	Shell string
	Args  []string
	Env   []string
	Cols  uint16
	Rows  uint16
	Dir   string
}

// New creates a new PTY with the given configuration.
func New(cfg Config) (*PTY, error) {
	if cfg.Cols == 0 {
		cfg.Cols = 120
	}
	if cfg.Rows == 0 {
		cfg.Rows = 40
	}
	if cfg.Shell == "" {
		cfg.Shell = defaultShell()
	}

	p := &PTY{
		wsCol: cfg.Cols,
		wsRow: cfg.Rows,
	}

	var err error
	if runtime.GOOS == "windows" {
		err = p.startWindows(cfg)
	} else {
		err = p.startUnix(cfg)
	}
	if err != nil {
		return nil, fmt.Errorf("pty start: %w", err)
	}

	return p, nil
}

// Write sends data to the PTY's stdin.
func (p *PTY) Write(data []byte) (int, error) {
	if p.stdin != nil {
		return p.stdin.Write(data)
	}
	if wp, ok := p.platform.(platformPTY); ok {
		return wp.Write(data)
	}
	return 0, fmt.Errorf("pty: no input available")
}

// Read reads output from the PTY.
func (p *PTY) Read(buf []byte) (int, error) {
	if p.stdout != nil {
		return p.stdout.Read(buf)
	}
	if wp, ok := p.platform.(platformPTY); ok {
		return wp.Read(buf)
	}
	return 0, io.EOF
}

// Resize changes the terminal dimensions.
func (p *PTY) Resize(cols, rows uint16) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.wsCol = cols
	p.wsRow = rows
	if wp, ok := p.platform.(platformPTY); ok {
		return wp.Resize(cols, rows)
	}
	return nil
}

// Close shuts down the PTY and child process.
func (p *PTY) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.running = false
	if p.stdin != nil {
		p.stdin.Close()
	}
	if p.cmd != nil && p.cmd.Process != nil {
		p.cmd.Process.Kill()
		p.cmd.Wait()
	}
	if wp, ok := p.platform.(platformPTY); ok {
		return wp.Close()
	}
	return nil
}

// Running returns whether the PTY process is alive.
func (p *PTY) Running() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.running
}

// Size returns the current terminal dimensions.
func (p *PTY) Size() (cols, rows uint16) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.wsCol, p.wsRow
}

// Process returns the underlying process, if any.
func (p *PTY) Process() *os.Process {
	if p.cmd != nil {
		return p.cmd.Process
	}
	return nil
}

func defaultShell() string {
	if runtime.GOOS == "windows" {
		return "cmd.exe"
	}
	if sh := os.Getenv("SHELL"); sh != "" {
		return sh
	}
	return "/bin/bash"
}

// ─── Unix PTY (pipe-based fallback) ──────────────────────────

func (p *PTY) startUnix(cfg Config) error {
	args := append([]string{}, cfg.Args...)
	p.cmd = exec.Command(cfg.Shell, args...)
	if cfg.Dir != "" {
		p.cmd.Dir = cfg.Dir
	}
	if len(cfg.Env) > 0 {
		p.cmd.Env = cfg.Env
	} else {
		p.cmd.Env = os.Environ()
	}

	stdin, err := p.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("stdin pipe: %w", err)
	}
	stdout, err := p.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := p.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderr pipe: %w", err)
	}

	p.stdout = io.MultiReader(stdout, stderr)
	p.stdin = stdin

	if err := p.cmd.Start(); err != nil {
		return fmt.Errorf("cmd start: %w", err)
	}
	p.running = true
	return nil
}
