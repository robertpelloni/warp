package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/robertpelloni/warp-rebuild/pkg/terminal"
	"golang.org/x/term"
)

func main() {
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
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if fd, isTerm := os.Stdout.Fd(), term.IsTerminal(int(os.Stdout.Fd())); isTerm {
				w, h, _ := term.GetSize(int(fd))
				session.Resize(w, h)
			}
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize

	// Set terminal to raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err == nil {
		defer term.Restore(int(os.Stdin.Fd()), oldState)
	}

	// Copy stdin to PTY and PTY to stdout
	go func() { io.Copy(session, os.Stdin) }()
	io.Copy(os.Stdout, session)
}
