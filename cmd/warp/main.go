package main

import (
	"fmt"
	"os"

	"github.com/robertpelloni/warp/pkg/app"
)

func main() {
	fmt.Println("Warp Go - Agentic Development Environment")
	fmt.Println("Port of warp-dev/warp from Rust to Go")

	cfg := app.Config{
		Shell: "",
		Theme: "",
	}

	// Parse command line flags
	for i, arg := range os.Args[1:] {
		_ = i
		switch arg {
		case "--terminal-only":
			cfg.TerminalOnly = true
		case "--editor-only":
			cfg.EditorOnly = true
		case "--theme":
			if i+1 < len(os.Args[1:]) {
				cfg.Theme = os.Args[i+2]
			}
		case "--shell":
			if i+1 < len(os.Args[1:]) {
				cfg.Shell = os.Args[i+2]
			}
		case "--help", "-h":
			fmt.Println(`
Warp Go - Agentic Development Environment

Usage: warp-go [options]

Options:
  --terminal-only    Terminal only mode (no editor pane)
  --editor-only      Editor only mode (no terminal pane)
  --theme <name>     Set theme (Standard, Dracula, Monokai, Nord, One Dark, Catppuccin)
  --shell <path>     Set shell (default: powershell or cmd.exe)
  --help, -h         Show this help

Keyboard Shortcuts:
  Ctrl+Enter    Execute command
  Ctrl+L        Clear output
  Ctrl+K        Command palette
  Ctrl+T        New session
  Ctrl+Shift+T  Cycle theme
  Ctrl+Up/Down  Command history

Built-in Commands:
  /help            Show Warp commands
  /ai <query>      Ask AI assistant
  /aliases         List aliases
  /alias <n> <cmd> Create alias

Built-in Aliases:
  ll  -> ls -la    gs -> git status    gd -> git diff
  gl  -> git log   gp -> git push      gpl -> git pull`)
			os.Exit(0)
		}
	}

	a := app.New(cfg)
	a.Run()
}
