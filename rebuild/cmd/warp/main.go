package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/robertpelloni/warp-rebuild/pkg/agent/service"
	"github.com/robertpelloni/warp-rebuild/pkg/terminal"
	"golang.org/x/term"
)

func main() {
	port := flag.Int("port", 0, "Port to run the agent service on (enables service mode)")
	flag.Parse()

	if *port != 0 {
		runService(*port)
		return
	}

	runTerminalHarness()
}

func runService(port int) {
	fmt.Printf("Warp Rebuild (Go) - Agent Service Mode on port %d\n", port)
	svc := service.NewAgentService(port)
	if err := svc.Start(context.Background()); err != nil {
		fmt.Printf("Service error: %v\n", err)
	}
}

func runTerminalHarness() {
	fmt.Println("Warp Rebuild (Go) - Terminal Harness Mode")

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "bash"
	}

	session, err := terminal.NewLocalSession(shell)
	if err != nil {
		fmt.Printf("Error creating session: %v\n", err)
		return
	}
	defer session.Close()

	// Handle terminal resizing
	setupResize(session)

	// Set terminal to raw mode
	if term.IsTerminal(int(os.Stdin.Fd())) {
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err == nil {
			defer term.Restore(int(os.Stdin.Fd()), oldState)
		}
	}

	// Copy stdin to PTY and PTY to stdout
	go func() { io.Copy(session, os.Stdin) }()
	io.Copy(os.Stdout, session)
}
