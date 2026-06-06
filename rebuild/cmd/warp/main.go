package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/robertpelloni/warp-rebuild/pkg/agent/service"
	"github.com/robertpelloni/warp-rebuild/pkg/terminal"
	"golang.org/x/term"
)

func main() {
	port := flag.Int("port", 0, "service port")
	useTUI := flag.Bool("tui", false, "run with experimental TUI")
	flag.Parse()

	if *port != 0 {
		s := service.NewAgentService(*port)
		s.Start(context.Background())
		return
	}

	if *useTUI {
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running TUI: %v", err)
			os.Exit(1)
		}
		return
	}

	runTerm()
}

func runTerminalHarness() { // For Unix signal handlers to bind to
}

func runTerm() {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "bash"
	}
	s, err := terminal.NewLocalSession(shell)
	if err != nil {
		fmt.Printf("Error creating local session: %v\n", err)
		os.Exit(1)
	}
	setupResize(s)
	if term.IsTerminal(int(os.Stdin.Fd())) {
		old, _ := term.MakeRaw(int(os.Stdin.Fd()))
		defer term.Restore(int(os.Stdin.Fd()), old)
	}
	go io.Copy(s, os.Stdin)
	io.Copy(os.Stdout, s)
}
