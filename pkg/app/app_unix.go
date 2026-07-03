//go:build !windows

package app

import (
	"context"
	"fmt"

	"github.com/robertpelloni/warp/pkg/shadowpilot"
)

type Config struct {
	Shell        string
	Theme        string
	TerminalOnly bool
	EditorOnly   bool
}

type WarpApp struct {
	cfg         Config
	shadowPilot *shadowpilot.Pilot
	ctx         context.Context
	cancel      context.CancelFunc
}

func New(cfg Config) *WarpApp {
	ctx, cancel := context.WithCancel(context.Background())
	return &WarpApp{
		cfg:    cfg,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (a *WarpApp) Run() {
	a.shadowPilot = shadowpilot.New(func(msg string) {
		fmt.Println(msg)
	})
	a.shadowPilot.Start(a.ctx)

	fmt.Println("Warp Go is currently a Windows-only application, as it depends on Win32 GUI APIs.")
}
