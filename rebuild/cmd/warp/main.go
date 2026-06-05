package main

import (
	"context"
	"flag"
	"io"
	"os"
	"github.com/robertpelloni/warp-rebuild/pkg/terminal"
	"github.com/robertpelloni/warp-rebuild/pkg/agent/service"
	"golang.org/x/term"
)

func main() {
	port := flag.Int("port", 0, "service port")
	flag.Parse()
	if *port != 0 {
		s := service.NewAgentService(*port)
		s.Start(context.Background())
		return
	}
	runTerm()
}

func runTerminalHarness() { // For Unix signal handlers to bind to
}

func runTerm() {
	shell := os.Getenv("SHELL")
	if shell == "" { shell = "bash" }
	s, _ := terminal.NewLocalSession(shell)
	setupResize(s)
	if term.IsTerminal(int(os.Stdin.Fd())) {
		old, _ := term.MakeRaw(int(os.Stdin.Fd()))
		defer term.Restore(int(os.Stdin.Fd()), old)
	}
	go io.Copy(s, os.Stdin)
	io.Copy(os.Stdout, s)
}
