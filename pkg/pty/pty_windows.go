//go:build windows

package pty

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"unsafe"

	"golang.org/x/sys/windows"
)

type winPTY struct {
	conPty     uintptr // HPCON
	inputPipe  uintptr // HANDLE - write side
	outputPipe uintptr // HANDLE - read side
	cmd        *exec.Cmd
}

func (p *PTY) startWindows(cfg Config) error {
	wp, err := newConPTY(cfg)
	if err != nil {
		// Fall back to pipe-based mode if ConPTY fails
		return p.startWindowsFallback(cfg)
	}
	p.platform = wp
	p.running = true
	return nil
}

func (p *PTY) startWindowsFallback(cfg Config) error {
	args := append([]string{}, cfg.Args...)
	p.cmd = exec.Command(cfg.Shell, args...)
	if cfg.Dir != "" {
		p.cmd.Dir = cfg.Dir
	}
	if len(cfg.Env) > 0 {
		p.cmd.Env = cfg.Env
	} else {
		p.cmd.Env = append(os.Environ(), "TERM=xterm-256color")
	}

	stdin, err := p.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("stdin pipe: %w", err)
	}
	stdout, err := p.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := p.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderr pipe: %w", err)
	}

	p.stdout = io.MultiReader(stdout, stderr)
	p.stdin = stdin

	if err := p.cmd.Start(); err != nil {
		return fmt.Errorf("cmd start: %w", err)
	}
	p.running = true
	return nil
}

func newConPTY(cfg Config) (*winPTY, error) {
	wp := &winPTY{}

	// Create pipes for ConPTY I/O
	var hInputRead, hInputWrite uintptr
	ret, _, _ := procCreatePipe.Call(
		uintptr(unsafe.Pointer(&hInputRead)),
		uintptr(unsafe.Pointer(&hInputWrite)),
		0, 0,
	)
	if ret == 0 {
		return nil, fmt.Errorf("CreatePipe(input) failed")
	}

	var hOutputRead, hOutputWrite uintptr
	ret, _, _ = procCreatePipe.Call(
		uintptr(unsafe.Pointer(&hOutputRead)),
		uintptr(unsafe.Pointer(&hOutputWrite)),
		0, 0,
	)
	if ret == 0 {
		procCloseHandle.Call(hInputRead)
		procCloseHandle.Call(hInputWrite)
		return nil, fmt.Errorf("CreatePipe(output) failed")
	}

	// Create pseudo console
	size := coord{X: int16(cfg.Cols), Y: int16(cfg.Rows)}
	var hPC uintptr
	ret, _, _ = procCreatePseudoConsole.Call(
		uintptr(unsafe.Pointer(&size)),
		hInputRead,
		hOutputWrite,
		0,
		uintptr(unsafe.Pointer(&hPC)),
	)
	if ret != 0 {
		procCloseHandle.Call(hInputRead)
		procCloseHandle.Call(hInputWrite)
		procCloseHandle.Call(hOutputRead)
		procCloseHandle.Call(hOutputWrite)
		return nil, fmt.Errorf("CreatePseudoConsole failed: 0x%x", ret)
	}

	// Close the ConPTY-owned ends
	procCloseHandle.Call(hInputRead)
	procCloseHandle.Call(hOutputWrite)

	// Start child process in ConPTY
	cmdline := cfg.Shell
	for _, arg := range cfg.Args {
		cmdline += " " + arg
	}

	if err := startConPTYProcess(hPC, cmdline, cfg.Dir, cfg.Env); err != nil {
		procClosePseudoConsole.Call(hPC)
		procCloseHandle.Call(hInputWrite)
		procCloseHandle.Call(hOutputRead)
		return nil, fmt.Errorf("start conpty process: %w", err)
	}

	wp.conPty = hPC
	wp.inputPipe = hInputWrite
	wp.outputPipe = hOutputRead
	return wp, nil
}

func (wp *winPTY) Write(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}
	var written uintptr
	ret, _, _ := procWriteFile.Call(
		wp.inputPipe,
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
		uintptr(unsafe.Pointer(&written)),
		0,
	)
	if ret == 0 {
		return 0, fmt.Errorf("WriteFile failed")
	}
	return int(written), nil
}

func (wp *winPTY) Read(buf []byte) (int, error) {
	// First check if data is available using PeekNamedPipe
	var bytesAvail uintptr
	ret, _, _ := procPeekNamedPipe.Call(
		wp.outputPipe,
		0,                                    // buffer
		0,                                    // buffer size
		0,                                    // bytes read
		uintptr(unsafe.Pointer(&bytesAvail)), // bytes available
		0,                                    // bytes left this message
	)
	if ret == 0 || bytesAvail == 0 {
		return 0, fmt.Errorf("no data available")
	}

	// Now read the available data
	var bytesRead uintptr
	ret, _, _ = procReadFile.Call(
		wp.outputPipe,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
		uintptr(unsafe.Pointer(&bytesRead)),
		0,
	)
	if ret == 0 {
		return 0, fmt.Errorf("ReadFile failed")
	}
	return int(bytesRead), nil
}

func (wp *winPTY) Resize(cols, rows uint16) error {
	if wp.conPty == 0 {
		return nil
	}
	size := coord{X: int16(cols), Y: int16(rows)}
	ret, _, _ := procResizePseudoConsole.Call(wp.conPty, uintptr(unsafe.Pointer(&size)))
	if ret != 0 {
		return fmt.Errorf("ResizePseudoConsole failed: 0x%x", ret)
	}
	return nil
}

func (wp *winPTY) Close() error {
	if wp.conPty != 0 {
		procClosePseudoConsole.Call(wp.conPty)
		wp.conPty = 0
	}
	if wp.inputPipe != 0 {
		procCloseHandle.Call(wp.inputPipe)
		wp.inputPipe = 0
	}
	if wp.outputPipe != 0 {
		procCloseHandle.Call(wp.outputPipe)
		wp.outputPipe = 0
	}
	return nil
}

func startConPTYProcess(hPC uintptr, cmdline, dir string, env []string) error {
	// Build STARTUPINFOEXW
	var si startupInfoExW
	si.StartupInfo.Cb = uint32(unsafe.Sizeof(si))

	// Calculate attribute list size
	var attrListSize uintptr
	procInitializeProcThreadAttributeList.Call(0, 1, 0, uintptr(unsafe.Pointer(&attrListSize)))

	attrBuf := make([]byte, attrListSize)
	si.AttrList = uintptr(unsafe.Pointer(&attrBuf[0]))

	ret, _, _ := procInitializeProcThreadAttributeList.Call(
		si.AttrList, 1, 0, uintptr(unsafe.Pointer(&attrListSize)),
	)
	if ret == 0 {
		return fmt.Errorf("InitializeProcThreadAttributeList failed")
	}
	defer procDeleteProcThreadAttributeList.Call(si.AttrList)

	// Set pseudo console attribute
	ret, _, _ = procUpdateProcThreadAttribute.Call(
		si.AttrList, 0,
		0x00020000, // PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE
		hPC, uintptr(unsafe.Sizeof(hPC)),
		0, 0,
	)
	if ret == 0 {
		return fmt.Errorf("UpdateProcThreadAttribute failed")
	}

	// Build command line UTF16
	cmdLineUTF16 := stringToUTF16Ptr(cmdline)

	var dirPtr *uint16
	if dir != "" {
		dirPtr = stringToUTF16Ptr(dir)
	}

	var pi processInformation
	ret, _, _ = procCreateProcessW.Call(
		0, // app name
		uintptr(unsafe.Pointer(cmdLineUTF16)),
		0, 0, // security
		0,                     // inherit handles
		0x00008000|0x00000040, // CREATE_UNICODE_ENVIRONMENT | EXTENDED_STARTUPINFO_PRESENT
		0,                     // environment
		uintptr(unsafe.Pointer(dirPtr)),
		uintptr(unsafe.Pointer(&si.StartupInfo)),
		uintptr(unsafe.Pointer(&pi)),
	)
	if ret == 0 {
		return fmt.Errorf("CreateProcessW failed")
	}

	procCloseHandle.Call(pi.Thread)
	procCloseHandle.Call(pi.Process)
	return nil
}

// ─── Windows types ───────────────────────────────────────────

type coord struct {
	X int16
	Y int16
}

type startupInfoExW struct {
	StartupInfo startupInfoW
	AttrList    uintptr
}

type startupInfoW struct {
	Cb            uint32
	Reserved1     *uint16
	Desktop       *uint16
	Title         *uint16
	X             int32
	Y             int32
	XSize         int32
	YSize         int32
	XCountChars   int32
	YCountChars   int32
	FillAttribute uint32
	Flags         uint32
	ShowWindow    uint16
	Reserved2     uint16
	Reserved3     *byte
	StdInput      uintptr
	StdOutput     uintptr
	StdError      uintptr
}

type processInformation struct {
	Process  uintptr
	Thread   uintptr
	ProcID   uint32
	ThreadID uint32
}

// ─── Windows proc table ─────────────────────────────────────

var (
	modkernel32 = windows.NewLazyDLL("kernel32.dll")

	procCreatePseudoConsole               = modkernel32.NewProc("CreatePseudoConsole")
	procResizePseudoConsole               = modkernel32.NewProc("ResizePseudoConsole")
	procClosePseudoConsole                = modkernel32.NewProc("ClosePseudoConsole")
	procCreateProcessW                    = modkernel32.NewProc("CreateProcessW")
	procCreatePipe                        = modkernel32.NewProc("CreatePipe")
	procWriteFile                         = modkernel32.NewProc("WriteFile")
	procReadFile                          = modkernel32.NewProc("ReadFile")
	procCloseHandle                       = modkernel32.NewProc("CloseHandle")
	procInitializeProcThreadAttributeList = modkernel32.NewProc("InitializeProcThreadAttributeList")
	procUpdateProcThreadAttribute         = modkernel32.NewProc("UpdateProcThreadAttribute")
	procDeleteProcThreadAttributeList     = modkernel32.NewProc("DeleteProcThreadAttributeList")
	procPeekNamedPipe                     = modkernel32.NewProc("PeekNamedPipe")
)

func stringToUTF16Ptr(s string) *uint16 {
	return windows.StringToUTF16Ptr(s)
}
