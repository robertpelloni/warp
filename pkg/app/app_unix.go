//go:build !windows

package app

import (
	"context"
	"fmt"
	"os"

	"github.com/robertpelloni/warp/pkg/agent"
	"github.com/robertpelloni/warp/pkg/shadowpilot"
)

type Config struct {
	Shell        string
	Theme        string
	TerminalOnly bool
	EditorOnly   bool
}

type WarpApp struct {
	cfg          Config
	shadowPilot  *shadowpilot.Pilot
	orchestrator *agent.Orchestrator
	ctx          context.Context
	cancel       context.CancelFunc
}

func New(cfg Config) *WarpApp {
	ctx, cancel := context.WithCancel(context.Background())
	app := &WarpApp{
		cfg:    cfg,
		ctx:    ctx,
		cancel: cancel,
	}

	cwd, _ := os.Getwd()
	if orch, err := agent.NewOrchestrator(cwd); err == nil {
		app.orchestrator = orch
		go func() {
			_ = app.orchestrator.ConnectRemote()
			fmt.Printf("[Agent Orchestrator] Status: %s\n", app.orchestrator.Status())
		}()
	}

	return app
}

func (a *WarpApp) Run() {
	a.shadowPilot = shadowpilot.New(func(msg string) {
		fmt.Println(msg)
	})
	a.shadowPilot.Start(a.ctx)

	fmt.Println("Warp Go is currently a Windows-only application, as it depends on Win32 GUI APIs.")
}
