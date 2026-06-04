//go:build !windows

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/robertpelloni/warp-rebuild/pkg/terminal"
	"golang.org/x/term"
)

func setupResize(session terminal.Session) {
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
}
