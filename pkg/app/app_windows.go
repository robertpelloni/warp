//go:build windows

package app

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/robertpelloni/warp/pkg/command"
	"github.com/robertpelloni/warp/pkg/editor"
	"github.com/robertpelloni/warp/pkg/session"
	"github.com/robertpelloni/warp/pkg/terminal"
	"github.com/robertpelloni/warp/pkg/theme"
	"github.com/robertpelloni/warp/pkg/win32"
)

// Win32 constants not in win32 package
const (
	CW_USEDEFAULT = ^0
	IDC_ARROW     = 32512
	CS_HREDRAW    = 0x0002
	CS_VREDRAW    = 0x0001
	WM_VSCROLL    = 0x0115

	// Child control IDs
	IDC_OUTPUT_EDIT  = 100
	IDC_INPUT_EDIT   = 101
	IDC_BLOCKS_LIST  = 102
	IDC_STATUS_BAR   = 103
	IDC_PROMPT_LABEL = 104
	IDC_EDITOR_EDIT  = 105

	// Tray identifier and custom message
	WM_APP_TRAYMSG = 0x8000 + 100
	ID_TRAY_ICON   = 1001

	// Tray Menu Command IDs
	ID_TRAY_SHOW = 2001
	ID_TRAY_HIDE = 2002
	ID_TRAY_QUIT = 2003
)

// WarpApp is the main application.
type WarpApp struct {
	mu sync.Mutex

	sessMgr    *session.Manager
	cmdEngine  *command.Engine
	editEngine *editor.Editor
	warpTheme  *theme.WarpTheme
	ctx        context.Context
	cancel     context.CancelFunc
	shell      string

	// Win32 window handles
	hMainWnd    win32.HWND
	hOutput     win32.HWND
	hInput      win32.HWND
	hBlocks     win32.HWND
	hEditor     win32.HWND
	hStatusBar  win32.HWND
	hPrompt     win32.HWND
	hFont       win32.HFONT
	hBoldFont   win32.HFONT
	hPromptFont win32.HFONT

	// System Tray
	trayData win32.NOTIFYICONDATAW

	// Colors (Win32 COLORREF)
	bgColor      uint32
	textColor    uint32
	accentColor  uint32
	promptColor  uint32
	surfaceColor uint32

	// State
	inputHistory   []string
	histIdx        int
	origInputProc  uintptr
	origEditorProc uintptr
}

// Config for the application.
type Config struct {
	Shell        string
	TerminalOnly bool
	EditorOnly   bool
	Theme        string
}

// New creates a new Warp application.
func New(cfg Config) *WarpApp {
	if cfg.Shell == "" {
		cfg.Shell = defaultShell()
	}
	ctx, cancel := context.WithCancel(context.Background())

	a := &WarpApp{
		sessMgr:    session.NewManager(),
		cmdEngine:  command.NewEngine(),
		editEngine: editor.New(),
		ctx:        ctx,
		cancel:     cancel,
		shell:      cfg.Shell,
	}

	if cfg.Theme != "" {
		a.warpTheme = theme.GetTheme(cfg.Theme)
	} else {
		a.warpTheme = theme.DefaultTheme()
	}
	a.updateColors()
	return a
}

func (a *WarpApp) updateColors() {
	a.bgColor = a.warpTheme.COLORREF(theme.ColorBackground)
	a.textColor = a.warpTheme.COLORREF(theme.ColorText)
	a.accentColor = a.warpTheme.COLORREF(theme.ColorAccent)
	a.promptColor = a.warpTheme.COLORREF(theme.ColorPrompt)
	a.surfaceColor = a.warpTheme.COLORREF(theme.ColorSurface)
}

// Run creates windows, initializes the tray, and runs the message loop.
func (a *WarpApp) Run() {
	hInstance := win32.GetModuleHandle()
	className := "WarpGoMainWnd"

	// Window procedure callback (must persist for lifetime)
	callback := syscall.NewCallback(func(hWnd uintptr, msg uint32, wParam uintptr, lParam uintptr) uintptr {
		return a.mainWndProc(win32.HWND(hWnd), msg, wParam, lParam)
	})

	cursor := win32.LoadCursor(0, uintptr(IDC_ARROW))
	bgBrush := win32.CreateSolidBrush(a.bgColor)
	classNamePtr := windows.StringToUTF16Ptr(className)

	wc := win32.WNDCLASSEXW{
		CbSize:        uint32(unsafe.Sizeof(win32.WNDCLASSEXW{})),
		Style:         CS_HREDRAW | CS_VREDRAW,
		LpfnWndProc:   callback,
		HInstance:     hInstance,
		HCursor:       cursor,
		HBrBackground: bgBrush,
		LpszClassName: classNamePtr,
	}
	win32.RegisterClassEx(&wc)

	// Create main window
	titlePtr := windows.StringToUTF16Ptr("Warp Go - Unified Agentic Development Dashboard")
	a.hMainWnd = win32.CreateWindowEx(
		win32.WS_EX_APPWINDOW,
		classNamePtr,
		titlePtr,
		win32.WS_OVERLAPPEDWINDOW,
		CW_USEDEFAULT, CW_USEDEFAULT, 1280, 860,
		0, 0, hInstance, 0,
	)

	// Fonts
	a.hFont = a.createMonoFont(15, win32.FW_NORMAL)
	a.hBoldFont = a.createMonoFont(15, win32.FW_BOLD)
	a.hPromptFont = a.createMonoFont(20, win32.FW_BOLD)

	a.createControls(hInstance)
	a.initTrayIcon(hInstance)
	a.initTerminal()

	win32.ShowWindow(a.hMainWnd, win32.SW_SHOWDEFAULT)
	win32.UpdateWindow(a.hMainWnd)

	// Refresh timer (50ms)
	win32.SetTimer(a.hMainWnd, 1, 50, 0)

	// Message loop
	var msg win32.MSG
	for {
		ret := win32.GetMessage(&msg, 0, 0, 0)
		if ret == 0 {
			break
		}
		win32.TranslateMessage(&msg)
		win32.DispatchMessage(&msg)
	}
	a.Quit()
}

func (a *WarpApp) initTrayIcon(hInstance win32.HINSTANCE) {
	hIcon := win32.LoadStandardIcon(win32.IDI_APPLICATION)
	a.trayData = win32.NOTIFYICONDATAW{
		CbSize:           uint32(unsafe.Sizeof(win32.NOTIFYICONDATAW{})),
		HWnd:             a.hMainWnd,
		UID:              ID_TRAY_ICON,
		UFlags:           win32.NIF_MESSAGE | win32.NIF_ICON | win32.NIF_TIP,
		UCallbackMessage: WM_APP_TRAYMSG,
		HIcon:            hIcon,
	}
	tip := "Warp Go - Unified Agentic Dashboard"
	tipUTF16 := windows.StringToUTF16(tip)
	copy(a.trayData.SzTip[:], tipUTF16)

	win32.ShellNotifyIcon(win32.NIM_ADD, &a.trayData)
}

func (a *WarpApp) removeTrayIcon() {
	if a.trayData.HWnd != 0 {
		win32.ShellNotifyIcon(win32.NIM_DELETE, &a.trayData)
		a.trayData.HWnd = 0
	}
}

func (a *WarpApp) showTrayMenu() {
	var pt win32.POINT
	win32.GetCursorPos(&pt)

	hMenu := win32.CreatePopupMenu()
	win32.AppendMenu(hMenu, win32.MF_STRING, ID_TRAY_SHOW, "Show Dashboard")
	win32.AppendMenu(hMenu, win32.MF_STRING, ID_TRAY_HIDE, "Minimize to Tray")
	win32.AppendMenu(hMenu, win32.MF_SEPARATOR, 0, "")
	win32.AppendMenu(hMenu, win32.MF_STRING, ID_TRAY_QUIT, "Quit Warp Servers & App")

	win32.SetForegroundWindow(a.hMainWnd)
	win32.TrackPopupMenu(hMenu, win32.TPM_RIGHTBUTTON|win32.TPM_BOTTOMALIGN, pt.X, pt.Y, a.hMainWnd)
	win32.DestroyMenu(hMenu)
}

func (a *WarpApp) createControls(hInstance win32.HINSTANCE) {
	// Status bar with detailed tooltips/hints
	a.hStatusBar = win32.CreateWindowEx(0,
		windows.StringToUTF16Ptr("STATIC"),
		windows.StringToUTF16Ptr(" [?] Dashboard Ready | Shell: "+a.shell+" | [Terminal: Ctrl+L Clear] | [Editor: /edit <path>, Ctrl+S Save] | [Commands: Ctrl+K Palette, /ai <query>]"),
		win32.WS_CHILD|win32.WS_VISIBLE,
		0, 0, 0, 0,
		a.hMainWnd, IDC_STATUS_BAR, hInstance, 0,
	)
	win32.SendMessage(a.hStatusBar, win32.WM_SETFONT, uintptr(a.hFont), 1)

	// Command blocks sidebar
	a.hBlocks = win32.CreateWindowEx(win32.WS_EX_CLIENTEDGE,
		windows.StringToUTF16Ptr("EDIT"),
		nil,
		win32.WS_CHILD|win32.WS_VISIBLE|win32.ES_MULTILINE|win32.ES_AUTOVSCROLL|win32.ES_READONLY|win32.WS_VSCROLL,
		0, 0, 0, 0,
		a.hMainWnd, IDC_BLOCKS_LIST, hInstance, 0,
	)
	win32.SendMessage(a.hBlocks, win32.WM_SETFONT, uintptr(a.hFont), 1)

	// Output area (Terminal)
	a.hOutput = win32.CreateWindowEx(win32.WS_EX_CLIENTEDGE,
		windows.StringToUTF16Ptr("EDIT"),
		nil,
		win32.WS_CHILD|win32.WS_VISIBLE|win32.ES_MULTILINE|win32.ES_AUTOVSCROLL|win32.ES_READONLY|win32.WS_VSCROLL,
		0, 0, 0, 0,
		a.hMainWnd, IDC_OUTPUT_EDIT, hInstance, 0,
	)
	win32.SendMessage(a.hOutput, win32.WM_SETFONT, uintptr(a.hFont), 1)

	// Editor area (Unified Dashboard IDE Pane)
	a.hEditor = win32.CreateWindowEx(win32.WS_EX_CLIENTEDGE,
		windows.StringToUTF16Ptr("EDIT"),
		windows.StringToUTF16Ptr("// [?] Warp IDE Editor Pane (Unified Dashboard)\r\n// - Type `/edit <filepath>` at prompt to load buffer\r\n// - Press `Ctrl+S` anywhere to save modifications\r\n"),
		win32.WS_CHILD|win32.WS_VISIBLE|win32.ES_MULTILINE|win32.ES_AUTOVSCROLL|win32.WS_VSCROLL,
		0, 0, 0, 0,
		a.hMainWnd, IDC_EDITOR_EDIT, hInstance, 0,
	)
	win32.SendMessage(a.hEditor, win32.WM_SETFONT, uintptr(a.hFont), 1)
	// Subclass editor pane to intercept Ctrl+S
	a.origEditorProc = win32.GetWindowLongPtr(a.hEditor, win32.GWLP_WNDPROC)
	editorCallback := syscall.NewCallback(func(hWnd uintptr, msg uint32, wParam uintptr, lParam uintptr) uintptr {
		if msg == win32.WM_KEYDOWN {
			ctrl := win32.IsKeyDown(win32.VK_CONTROL)
			if ctrl && wParam == 'S' {
				a.saveActiveEditorBuffer()
				return 0
			}
		}
		return win32.CallWindowProc(a.origEditorProc, win32.HWND(hWnd), msg, wParam, lParam)
	})
	win32.SetWindowLongPtr(a.hEditor, win32.GWLP_WNDPROC, editorCallback)

	// Prompt label
	a.hPrompt = win32.CreateWindowEx(0,
		windows.StringToUTF16Ptr("STATIC"),
		windows.StringToUTF16Ptr(">"),
		win32.WS_CHILD|win32.WS_VISIBLE,
		0, 0, 0, 0,
		a.hMainWnd, IDC_PROMPT_LABEL, hInstance, 0,
	)
	win32.SendMessage(a.hPrompt, win32.WM_SETFONT, uintptr(a.hPromptFont), 1)

	// Input edit
	a.hInput = win32.CreateWindowEx(win32.WS_EX_CLIENTEDGE,
		windows.StringToUTF16Ptr("EDIT"),
		nil,
		win32.WS_CHILD|win32.WS_VISIBLE|win32.ES_AUTOHSCROLL,
		0, 0, 0, 0,
		a.hMainWnd, IDC_INPUT_EDIT, hInstance, 0,
	)
	win32.SendMessage(a.hInput, win32.WM_SETFONT, uintptr(a.hFont), 1)

	// Subclass the input field to intercept keyboard events
	a.origInputProc = win32.GetWindowLongPtr(a.hInput, win32.GWLP_WNDPROC)
	inputCallback := syscall.NewCallback(func(hWnd uintptr, msg uint32, wParam uintptr, lParam uintptr) uintptr {
		if msg == win32.WM_KEYDOWN && wParam == win32.VK_RETURN {
			a.executeInput()
			return 0
		}
		if msg == win32.WM_CHAR && wParam == 13 {
			return 0 // suppress beep
		}
		if msg == win32.WM_KEYDOWN {
			ctrl := win32.IsKeyDown(win32.VK_CONTROL)
			if ctrl && wParam == win32.VK_UP {
				a.historyUp()
				return 0
			}
			if ctrl && wParam == win32.VK_DOWN {
				a.historyDown()
				return 0
			}
			if ctrl && wParam == 'L' {
				a.clearOutput()
				return 0
			}
			if ctrl && wParam == 'K' {
				a.showCommandPalette()
				return 0
			}
		}
		return win32.CallWindowProc(a.origInputProc, win32.HWND(hWnd), msg, wParam, lParam)
	})
	win32.SetWindowLongPtr(a.hInput, win32.GWLP_WNDPROC, inputCallback)

	win32.SetFocus(a.hInput)

	a.appendOutput(a.welcomeMessage())
}

func (a *WarpApp) initTerminal() {
	sess, err := a.sessMgr.Create("Terminal 1", a.shell, 120, 40, func() {
		if a.hMainWnd != 0 {
			win32.PostMessage(a.hMainWnd, win32.WM_TIMER, 1, 0)
		}
	})
	if err != nil {
		a.appendOutput(fmt.Sprintf("\r\n[error] Failed to create terminal: %v\r\n", err))
		return
	}
	if err := sess.Terminal.Start(a.ctx); err != nil {
		a.appendOutput("\r\n[info] Terminal PTY not available, using shell exec mode\r\n")
	} else {
		a.appendOutput(fmt.Sprintf("\r\n[info] Terminal started: %s\r\n", a.shell))
	}
}

func (a *WarpApp) welcomeMessage() string {
	return "Warp Go - Unified Agentic Development Dashboard\r\n" +
		"Single-Dashboard Architecture with IDE Editor, Command Blocks, and Real PTY\r\n\r\n" +
		"Features:\r\n" +
		"  - Real PTY with ConPTY support\r\n" +
		"  - Command blocks (Warp-style output grouping)\r\n" +
		"  - Unified IDE Editor pane (/edit <file>)\r\n" +
		"  - System tray icon with background server control\r\n" +
		"  - AI integration (/ai <query>)\r\n" +
		"  - Command palette (Ctrl+K)\r\n" +
		"  - 7 built-in themes (Ctrl+Shift+T)\r\n" +
		"  - Shell aliases and auto-complete\r\n\r\n" +
		"Keyboard Shortcuts:\r\n" +
		"  Enter          Execute command\r\n" +
		"  Ctrl+L         Clear output\r\n" +
		"  Ctrl+K         Command palette\r\n" +
		"  Ctrl+T         New session\r\n" +
		"  Ctrl+Shift+T   Cycle theme\r\n" +
		"  Ctrl+Up/Down   Command history\r\n\r\n" +
		"Type /help for Warp commands, /edit <path> to open a file, or any shell command.\r\n\r\n"
}

// -- Main Window Procedure --

func (a *WarpApp) mainWndProc(hWnd win32.HWND, msg uint32, wParam uintptr, lParam uintptr) uintptr {
	switch msg {
	case win32.WM_SIZE:
		a.layoutControls()
		return 0

	case win32.WM_ERASEBKGND:
		hdc := win32.HDC(wParam)
		rc := win32.GetClientRect(hWnd)
		brush := win32.CreateSolidBrush(a.bgColor)
		win32.FillRect(hdc, &rc, brush)
		win32.DeleteObject(uintptr(brush))
		return 1

	case win32.WM_CTLCOLORSTATIC, win32.WM_CTLCOLOREDIT:
		hdc := win32.HDC(wParam)
		ctrlWnd := win32.HWND(lParam)
		win32.SetBkMode(hdc, win32.OPAQUE)

		switch ctrlWnd {
		case a.hInput:
			win32.SetBkColor(hdc, a.surfaceColor)
			win32.SetTextColor(hdc, a.textColor)
			return uintptr(win32.CreateSolidBrush(a.surfaceColor))
		case a.hStatusBar:
			win32.SetBkColor(hdc, a.surfaceColor)
			win32.SetTextColor(hdc, a.accentColor)
			return uintptr(win32.CreateSolidBrush(a.surfaceColor))
		case a.hPrompt:
			win32.SetBkColor(hdc, a.bgColor)
			win32.SetTextColor(hdc, a.promptColor)
			return uintptr(win32.CreateSolidBrush(a.bgColor))
		case a.hBlocks, a.hEditor:
			win32.SetBkColor(hdc, a.surfaceColor)
			win32.SetTextColor(hdc, a.textColor)
			return uintptr(win32.CreateSolidBrush(a.surfaceColor))
		default:
			win32.SetBkColor(hdc, a.bgColor)
			win32.SetTextColor(hdc, a.textColor)
			return uintptr(win32.CreateSolidBrush(a.bgColor))
		}

	case win32.WM_COMMAND:
		switch wParam {
		case ID_TRAY_SHOW:
			win32.ShowWindow(a.hMainWnd, win32.SW_RESTORE)
			win32.SetForegroundWindow(a.hMainWnd)
			return 0
		case ID_TRAY_HIDE:
			win32.ShowWindow(a.hMainWnd, win32.SW_HIDE)
			return 0
		case ID_TRAY_QUIT:
			a.Quit()
			win32.PostQuitMessage(0)
			return 0
		}

	case WM_APP_TRAYMSG:
		switch lParam {
		case win32.WM_RBUTTONUP:
			a.showTrayMenu()
		case win32.WM_LBUTTONDBLCLK:
			win32.ShowWindow(a.hMainWnd, win32.SW_RESTORE)
			win32.SetForegroundWindow(a.hMainWnd)
		}
		return 0

	case win32.WM_KEYDOWN:
		vk := wParam
		ctrl := win32.IsKeyDown(win32.VK_CONTROL)
		shift := win32.IsKeyDown(win32.VK_SHIFT)

		switch {
		case ctrl && vk == win32.VK_RETURN:
			a.executeInput()
			return 0
		case ctrl && vk == 'L':
			a.clearOutput()
			return 0
		case ctrl && vk == 'K':
			a.showCommandPalette()
			return 0
		case ctrl && shift && vk == 'T':
			a.cycleTheme()
			return 0
		case ctrl && vk == 'T':
			a.createNewSession()
			return 0
		case ctrl && vk == win32.VK_UP:
			a.historyUp()
			return 0
		case ctrl && vk == win32.VK_DOWN:
			a.historyDown()
			return 0
		}

	case win32.WM_CHAR:
		if wParam == 13 { // Enter
			a.executeInput()
			return 0
		}
		if wParam == 12 { // Ctrl+L
			a.clearOutput()
			return 0
		}
		if wParam == 11 { // Ctrl+K
			a.showCommandPalette()
			return 0
		}

	case win32.WM_TIMER:
		a.refreshTerminalOutput()
		return 0

	case win32.WM_DESTROY:
		win32.KillTimer(hWnd, 1)
		a.removeTrayIcon()
		win32.PostQuitMessage(0)
		return 0
	}

	return win32.DefWindowProc(hWnd, msg, wParam, lParam)
}

// -- Layout: Unified Single-Dashboard Layout --

func (a *WarpApp) layoutControls() {
	rc := win32.GetClientRect(a.hMainWnd)
	w := rc.Right - rc.Left
	h := rc.Bottom - rc.Top

	statusH := int32(28)
	inputH := int32(36)
	promptW := int32(30)
	pad := int32(4)
	blockW := int32(230)

	win32.MoveWindow(a.hStatusBar, 0, h-statusH, w, statusH, true)

	inputY := h - statusH - inputH
	win32.MoveWindow(a.hPrompt, pad, inputY+pad, promptW, inputH-pad*2, true)
	win32.MoveWindow(a.hInput, pad+promptW, inputY, w-pad*2-promptW, inputH, true)

	contentH := inputY - pad

	// Dashboard division:
	// Left: Terminal Output (50% of main area)
	// Center: Editor pane (50% of main area)
	// Right: Command Blocks sidebar (blockW)
	remainingW := w - blockW - pad*4
	if remainingW < 400 {
		// Responsive collapse if viewport is narrow
		blockW = 0
		remainingW = w - pad*2
	}

	terminalW := remainingW / 2
	editorW := remainingW - terminalW

	win32.MoveWindow(a.hOutput, pad, pad, terminalW, contentH, true)
	win32.MoveWindow(a.hEditor, pad*2+terminalW, pad, editorW, contentH, true)

	if blockW > 0 {
		win32.MoveWindow(a.hBlocks, pad*3+terminalW+editorW, pad, blockW, contentH, true)
	}
}

// -- Terminal & Editor Operations --

func (a *WarpApp) executeInput() {
	text := win32.GetWindowText(a.hInput)
	if text == "" {
		return
	}
	a.executeCommand(text)
	win32.SetWindowText(a.hInput, "")
	win32.SetFocus(a.hInput)
}

func (a *WarpApp) executeCommand(cmd string) {
	a.inputHistory = append(a.inputHistory, cmd)
	a.histIdx = len(a.inputHistory)

	cmdType := a.cmdEngine.Classify(cmd)
	expanded := a.cmdEngine.Expand(cmd)

	switch cmdType {
	case command.CmdWarp, command.CmdAI:
		result, _, err := a.cmdEngine.Execute(cmd)
		if err != nil {
			a.appendOutput(fmt.Sprintf("\r\nError: %v\r\n", err))
		} else {
			a.appendOutput(fmt.Sprintf("\r\n%s\r\n", result))
		}
	case command.CmdEditor:
		parts := strings.Fields(expanded)
		if len(parts) >= 2 {
			a.openEditor(parts[1])
		}
	case command.CmdAlias:
		a.appendOutput(fmt.Sprintf("\r\n> %s -> %s\r\n", cmd, expanded))
		a.sendToTerminal(expanded)
	default:
		a.appendOutput(fmt.Sprintf("\r\n> %s\r\n", cmd))
		a.sendToTerminal(expanded)
	}

	a.updateBlocksView()
	a.setStatus(fmt.Sprintf("Executed: %s | %s", cmd, time.Now().Format("15:04:05")))
}

func (a *WarpApp) saveActiveEditorBuffer() {
	buf := a.editEngine.Active()
	if buf == nil || buf.Path == "" {
		a.setStatus("No active file in editor to save")
		return
	}
	text := win32.GetWindowText(a.hEditor)
	// Strip header if present
	lines := strings.Split(text, "\r\n")
	var contentLines []string
	headerSkipped := false
	for _, line := range lines {
		if !headerSkipped && strings.HasPrefix(line, "// File: ") {
			headerSkipped = true
			continue
		}
		contentLines = append(contentLines, line)
	}
	actualContent := strings.Join(contentLines, "\n")
	a.editEngine.SetContent(buf.Path, actualContent)
	if err := a.editEngine.SaveToFile(buf.Path); err != nil {
		a.appendOutput(fmt.Sprintf("\r\n[Editor] Failed to save %s: %v\r\n", buf.Path, err))
		a.setStatus(fmt.Sprintf("Save failed: %v", err))
		return
	}
	a.appendOutput(fmt.Sprintf("\r\n[Editor] Successfully saved %s\r\n", buf.Path))
	a.setStatus(fmt.Sprintf("Saved: %s at %s", filepath.Base(buf.Path), time.Now().Format("15:04:05")))
}

func (a *WarpApp) openEditor(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		a.appendOutput(fmt.Sprintf("\r\n[Editor] Failed to open %s: %v\r\n", filePath, err))
		return
	}
	buf := a.editEngine.Open(filePath)
	buf.Content = string(content)
	buf.Language = editor.DetectLanguage(filePath)

	editorText := fmt.Sprintf("// File: %s (%s, %d lines)\r\n%s",
		filepath.Base(filePath), buf.Language, buf.LineCount(), strings.ReplaceAll(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n", "\r\n"))
	win32.SetWindowText(a.hEditor, editorText)
	a.appendOutput(fmt.Sprintf("\r\n[Editor] Loaded %s into dashboard editor\r\n", filePath))
}

func (a *WarpApp) sendToTerminal(cmd string) {
	sess := a.sessMgr.Active()
	if sess != nil && sess.Terminal != nil && sess.Terminal.Running() {
		sess.Terminal.SendCommand(cmd)
		return
	}
	a.executeShellCommand(cmd)
}

func (a *WarpApp) executeShellCommand(cmd string) {
	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		c = exec.Command("cmd", "/c", cmd)
	} else {
		c = exec.Command("sh", "-c", cmd)
	}
	c.Env = os.Environ()
	c.Dir, _ = os.Getwd()

	output, err := c.CombinedOutput()
	if err != nil {
		a.appendOutput(fmt.Sprintf("Error: %v\r\n", err))
	}
	if len(output) > 0 {
		a.appendOutput(string(output))
	}
}

func (a *WarpApp) appendOutput(text string) {
	win32.SendMessage(a.hOutput, win32.EM_SETSEL, ^uintptr(0), ^uintptr(0))
	win32.SendMessage(a.hOutput, win32.EM_REPLACESEL, 0, uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(text))))
	win32.SendMessage(a.hOutput, WM_VSCROLL, uintptr(win32.SB_BOTTOM), 0)
}

func (a *WarpApp) clearOutput() {
	win32.SetWindowText(a.hOutput, "")
	a.setStatus("Output cleared")
}

func (a *WarpApp) refreshTerminalOutput() {
	sess := a.sessMgr.Active()
	if sess == nil || sess.Terminal == nil || !sess.Terminal.Running() {
		return
	}
	output := sess.Terminal.GetOutput()
	if output == "" {
		return
	}
	cleanText, _ := terminal.ParseANSI(output)
	cleanText = strings.ReplaceAll(cleanText, "\n", "\r\n")
	win32.SetWindowText(a.hOutput, cleanText)
	win32.SendMessage(a.hOutput, WM_VSCROLL, uintptr(win32.SB_BOTTOM), 0)
}

func (a *WarpApp) updateBlocksView() {
	sess := a.sessMgr.Active()
	if sess == nil || sess.Terminal == nil {
		return
	}
	blocks := sess.Terminal.Blocks()
	var b strings.Builder
	b.WriteString("--- [?] Command Blocks (Execution History) ---\r\n\r\n")
	for i, block := range blocks {
		icon := ">"
		status := "*"
		switch block.Type {
		case terminal.BlockCommand:
			icon = ">"
		case terminal.BlockAI:
			icon = "*"
		case terminal.BlockError:
			icon = "x"
		}
		if block.Done {
			status = "+"
		}
		b.WriteString(fmt.Sprintf(" %s %s %s\r\n", icon, status, block.Command))
		if block.Done && block.Duration > 0 {
			b.WriteString(fmt.Sprintf("    %s\r\n", block.Duration.Round(time.Millisecond)))
		}
		if i < len(blocks)-1 {
			b.WriteString("\r\n")
		}
	}
	if len(blocks) == 0 {
		b.WriteString("  No commands yet.\r\n  Type below to start.\r\n")
	}
	win32.SetWindowText(a.hBlocks, b.String())
}

func (a *WarpApp) historyUp() {
	if a.histIdx > 0 {
		a.histIdx--
	}
	if a.histIdx < len(a.inputHistory) {
		win32.SetWindowText(a.hInput, a.inputHistory[a.histIdx])
	}
}

func (a *WarpApp) historyDown() {
	if a.histIdx < len(a.inputHistory)-1 {
		a.histIdx++
		win32.SetWindowText(a.hInput, a.inputHistory[a.histIdx])
	} else {
		a.histIdx = len(a.inputHistory)
		win32.SetWindowText(a.hInput, "")
	}
}

func (a *WarpApp) createNewSession() {
	count := a.sessMgr.Count() + 1
	name := fmt.Sprintf("Terminal %d", count)
	sess, err := a.sessMgr.Create(name, a.shell, 120, 40, func() {
		if a.hMainWnd != 0 {
			win32.PostMessage(a.hMainWnd, win32.WM_TIMER, 1, 0)
		}
	})
	if err != nil {
		a.appendOutput(fmt.Sprintf("\r\n[error] %v\r\n", err))
		return
	}
	if err := sess.Terminal.Start(a.ctx); err != nil {
		a.appendOutput("\r\n[info] Shell exec mode\r\n")
	}
	win32.SetWindowText(a.hOutput, "")
	a.appendOutput(fmt.Sprintf("Warp Go - %s - Ready\r\n\r\n", name))
	a.setStatus(fmt.Sprintf("Opened %s | Sessions: %d", name, a.sessMgr.Count()))
}

func (a *WarpApp) showCommandPalette() {
	var entries []string
	entries = append(entries,
		"/help - Show Warp commands",
		"/ai <query> - Ask AI assistant",
		"/edit <file> - Open file in Editor",
		"/aliases - List aliases",
		"New Tab (Ctrl+T)",
		"Clear Output (Ctrl+L)",
		"Switch Theme (Ctrl+Shift+T)",
	)
	for name, target := range a.cmdEngine.Aliases() {
		entries = append(entries, fmt.Sprintf("%s -> %s", name, target))
	}

	var b strings.Builder
	b.WriteString("\r\n--- Command Palette ---\r\n\r\n")
	for _, e := range entries {
		b.WriteString("  " + e + "\r\n")
	}
	b.WriteString("\r\nType a command or /help\r\n")
	a.appendOutput(b.String())
	win32.SetFocus(a.hInput)
}

func (a *WarpApp) cycleTheme() {
	themes := theme.ThemeNames()
	currentIdx := 0
	for i, name := range themes {
		if name == a.warpTheme.Name() {
			currentIdx = i
			break
		}
	}
	nextIdx := (currentIdx + 1) % len(themes)
	a.warpTheme = theme.GetTheme(themes[nextIdx])
	a.updateColors()
	a.setStatus(fmt.Sprintf("Theme: %s", themes[nextIdx]))
	win32.InvalidateRect(a.hMainWnd, nil, true)
}

func (a *WarpApp) setStatus(text string) {
	win32.SetWindowText(a.hStatusBar, "  "+text)
}

func (a *WarpApp) createMonoFont(size int32, weight int32) win32.HFONT {
	return win32.CreateFont(
		-size, 0,
		0, 0,
		weight,
		0, 0, 0,
		win32.DEFAULT_CHARSET,
		win32.OUT_DEFAULT_PRECIS,
		win32.CLIP_DEFAULT_PRECIS,
		win32.DEFAULT_QUALITY,
		win32.FIXED_PITCH|win32.FF_MODERN,
		"Cascadia Mono",
	)
}

// Quit shuts down the application and background servers.
func (a *WarpApp) Quit() {
	a.removeTrayIcon()
	a.cancel()
	for _, s := range a.sessMgr.List() {
		if s.Terminal != nil {
			s.Terminal.Stop()
		}
	}
	a.editEngine.Close()
}

func defaultShell() string {
	if runtime.GOOS == "windows" {
		if _, err := exec.LookPath("powershell"); err == nil {
			return "powershell"
		}
		return "cmd.exe"
	}
	if sh := os.Getenv("SHELL"); sh != "" {
		return sh
	}
	return "/bin/bash"
}
