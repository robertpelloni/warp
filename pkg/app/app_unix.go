//go:build !windows

package app

import "fmt"

type Config struct {
	Shell        string
	Theme        string
	TerminalOnly bool
	EditorOnly   bool
}

type WarpApp struct {
	cfg Config
}

func New(cfg Config) *WarpApp {
	return &WarpApp{cfg: cfg}
}

func (a *WarpApp) Run() {
	fmt.Println("Warp Go is currently a Windows-only application, as it depends on Win32 GUI APIs.")
}
