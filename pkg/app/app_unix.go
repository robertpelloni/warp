//go:build !windows

package app

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/robertpelloni/warp/pkg/command"
	"github.com/robertpelloni/warp/pkg/editor"
	"github.com/robertpelloni/warp/pkg/session"
	"github.com/robertpelloni/warp/pkg/theme"
)

// WarpApp is the main application fallback for non-Windows platforms.
type WarpApp struct {
	mu         sync.Mutex
	sessMgr    *session.Manager
	cmdEngine  *command.Engine
	editEngine *editor.Editor
	warpTheme  *theme.WarpTheme
	ctx        context.Context
	cancel     context.CancelFunc
	shell      string
}

// Config for the application.
type Config struct {
	Shell        string
	TerminalOnly bool
	EditorOnly   bool
	Theme        string
}

// New creates a new Warp application on non-Windows platforms.
func New(cfg Config) *WarpApp {
	if cfg.Shell == "" {
		cfg.Shell = defaultShell()
	}
	ctx, cancel := context.WithCancel(context.Background())

	a := &WarpApp{
		sessMgr:    session.NewManager(),
		cmdEngine:  command.NewEngine(),
		editEngine: editor.New(),
		ctx:        ctx,
		cancel:     cancel,
		shell:      cfg.Shell,
	}

	if cfg.Theme != "" {
		a.warpTheme = theme.GetTheme(cfg.Theme)
	} else {
		a.warpTheme = theme.DefaultTheme()
	}
	return a
}

// Run provides a headless or console fallback run loop for non-Windows environments.
func (a *WarpApp) Run() {
	fmt.Printf("[Warp Go] Running on non-Windows OS with shell: %s\n", a.shell)
	fmt.Println("[Warp Go] Native Win32 GUI is disabled. CLI / headless mode active.")
}

// Quit shuts down the application.
func (a *WarpApp) Quit() {
	a.cancel()
	for _, s := range a.sessMgr.List() {
		if s.Terminal != nil {
			s.Terminal.Stop()
		}
	}
	a.editEngine.Close()
}

func defaultShell() string {
	if sh := os.Getenv("SHELL"); sh != "" {
		return sh
	}
	if _, err := exec.LookPath("bash"); err == nil {
		return "bash"
	}
	return "sh"
}
