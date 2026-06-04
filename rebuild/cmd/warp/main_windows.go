// +build windows

package main

import (
	"github.com/robertpelloni/warp-rebuild/pkg/terminal"
)

func setupResize(session terminal.Session) {
	// SIGWINCH and raw PTY resizing are different on Windows
	// and are often handled by the terminal emulator/conpty itself.
}
