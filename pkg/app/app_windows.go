//go:build windows

package app

import (
	"context"
	"fmt"
	"os"
	"os/exec"
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
	IDC_TAB_CONTROL  = 106
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
	hTabs       win32.HWND
	hTooltip    win32.HWND
	hStatusBar  win32.HWND
	hPrompt     win32.HWND
	hFont       win32.HFONT
	hBoldFont   win32.HFONT
	hPromptFont win32.HFONT

	// Colors (Win32 COLORREF)
	bgColor      uint32
	textColor    uint32
	accentColor  uint32
	promptColor  uint32
	surfaceColor uint32

	// State
	inputHistory  []string
	histIdx       int
	origInputProc uintptr
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

// Run creates windows and runs the message loop.
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
	titlePtr := windows.StringToUTF16Ptr("Warp Go - Agentic Development Environment")
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

func (a *WarpApp) createControls(hInstance win32.HINSTANCE) {
	// Status bar
	a.hTabs = win32.CreateWindowEx(0,
		windows.StringToUTF16Ptr("SysTabControl32"),
		nil,
		win32.WS_CHILD|win32.WS_VISIBLE,
		0, 0, 0, 0,
		a.hMainWnd, IDC_TAB_CONTROL, hInstance, 0,
	)
	win32.SendMessage(a.hTabs, win32.WM_SETFONT, uintptr(a.hFont), 1)

	a.hStatusBar = win32.CreateWindowEx(0,
		windows.StringToUTF16Ptr("STATIC"),
		windows.StringToUTF16Ptr("  Ready | Shell: "+a.shell+" | Ctrl+K: Commands | Ctrl+L: Clear | Ctrl+T: New Tab | Ctrl+Shift+T: Theme"),
		win32.WS_CHILD|win32.WS_VISIBLE,
		0, 0, 0, 0,
		a.hMainWnd, IDC_STATUS_BAR, hInstance, 0,
	)
	win32.SendMessage(a.hStatusBar, win32.WM_SETFONT, uintptr(a.hFont), 1)

	// Blocks sidebar
	a.hBlocks = win32.CreateWindowEx(win32.WS_EX_CLIENTEDGE,
		windows.StringToUTF16Ptr("EDIT"),
		nil,
		win32.WS_CHILD|win32.WS_VISIBLE|win32.ES_MULTILINE|win32.ES_AUTOVSCROLL|win32.ES_READONLY|win32.WS_VSCROLL,
		0, 0, 0, 0,
		a.hMainWnd, IDC_BLOCKS_LIST, hInstance, 0,
	)
	win32.SendMessage(a.hBlocks, win32.WM_SETFONT, uintptr(a.hFont), 1)

	// Output area
	a.hOutput = win32.CreateWindowEx(win32.WS_EX_CLIENTEDGE,
		windows.StringToUTF16Ptr("EDIT"),
		nil,
		win32.WS_CHILD|win32.WS_VISIBLE|win32.ES_MULTILINE|win32.ES_AUTOVSCROLL|win32.ES_READONLY|win32.WS_VSCROLL,
		0, 0, 0, 0,
		a.hMainWnd, IDC_OUTPUT_EDIT, hInstance, 0,
	)
	win32.SendMessage(a.hOutput, win32.WM_SETFONT, uintptr(a.hFont), 1)

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

	// Tooltips
	a.hTooltip = win32.CreateWindowEx(0,
		windows.StringToUTF16Ptr(win32.TOOLTIPS_CLASSW),
		nil,
		win32.WS_POPUP|win32.WS_EX_TOPMOST,
		0, 0, 0, 0,
		a.hMainWnd, 0, hInstance, 0,
	)

	addTooltip := func(hwnd win32.HWND, text string) {
		ti := win32.TOOLINFO{
			CbSize:   uint32(unsafe.Sizeof(win32.TOOLINFO{})),
			UFlags:   win32.TTF_SUBCLASS | win32.TTF_IDISHWND,
			Hwnd:     uintptr(a.hMainWnd),
			UId:      uintptr(hwnd),
			LpszText: windows.StringToUTF16Ptr(text),
		}
		win32.SendMessage(a.hTooltip, win32.TTM_ADDTOOLW, 0, uintptr(unsafe.Pointer(&ti)))
	}

	addTooltip(a.hPrompt, "Terminal Input Prompt")
	addTooltip(a.hInput, "Type commands here. Use /help for Warp commands.")
	addTooltip(a.hBlocks, "Command history and execution blocks.")
	addTooltip(a.hTabs, "Active terminal sessions.")

	// Subclass the input field to intercept Enter key
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
	return "Warp Go - Agentic Development Environment\r\n" +
		"Port of warp-dev/warp from Rust to Go\r\n\r\n" +
		"Features:\r\n" +
		"  - Real PTY with ConPTY support\r\n" +
		"  - Command blocks (Warp-style output grouping)\r\n" +
		"  - Multiple sessions and tabs\r\n" +
		"  - AI integration (/ai <query>)\r\n" +
		"  - Command palette (Ctrl+K)\r\n" +
		"  - 7 built-in themes (Ctrl+Shift+T)\r\n" +
		"  - Shell aliases and auto-complete\r\n\r\n" +
		"Keyboard Shortcuts:\r\n" +
		"  Ctrl+Enter    Execute command\r\n" +
		"  Ctrl+L        Clear output\r\n" +
		"  Ctrl+K        Command palette\r\n" +
		"  Ctrl+T        New session\r\n" +
		"  Ctrl+Shift+T  Cycle theme\r\n" +
		"  Ctrl+Up/Down  Command history\r\n\r\n" +
		"Type /help for Warp commands, or any shell command below.\r\n\r\n"
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
		case a.hBlocks:
			win32.SetBkColor(hdc, a.surfaceColor)
			win32.SetTextColor(hdc, a.textColor)
			return uintptr(win32.CreateSolidBrush(a.surfaceColor))
		default:
			win32.SetBkColor(hdc, a.bgColor)
			win32.SetTextColor(hdc, a.textColor)
			return uintptr(win32.CreateSolidBrush(a.bgColor))
		}

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
		win32.PostQuitMessage(0)
		return 0
	}

	return win32.DefWindowProc(hWnd, msg, wParam, lParam)
}

// -- Layout --

func (a *WarpApp) layoutControls() {
	rc := win32.GetClientRect(a.hMainWnd)
	w := rc.Right - rc.Left
	h := rc.Bottom - rc.Top

	statusH := int32(28)
	inputH := int32(36)
	promptW := int32(30)
	pad := int32(4)
	tabH := int32(30)
	blockW := int32(250)

	win32.MoveWindow(a.hTabs, 0, 0, w, tabH, true)
	win32.MoveWindow(a.hStatusBar, 0, h-statusH, w, statusH, true)

	inputY := h - statusH - inputH
	win32.MoveWindow(a.hPrompt, pad, inputY+pad, promptW, inputH-pad*2, true)
	win32.MoveWindow(a.hInput, pad+promptW, inputY, w-pad*2-promptW, inputH, true)

	contentH := inputY - tabH

	outputW := w - blockW - pad*3
	editorW := w / 3
	if w < 1000 {
		editorW = 0
	}

	outputW = outputW - editorW
	if outputW < 200 {
		outputW = w - pad*2
		blockW = 0
	}

	win32.MoveWindow(a.hOutput, pad, pad+tabH, outputW, contentH-pad*2, true)

	if blockW > 0 {
		win32.MoveWindow(a.hBlocks, outputW+pad*2, pad+tabH, blockW-pad, contentH-pad*2, true)
		win32.ShowWindow(a.hBlocks, win32.SW_SHOW)
	} else {
		win32.ShowWindow(a.hBlocks, win32.SW_HIDE)
	}

	if editorW > 0 {
		win32.MoveWindow(a.hEditor, outputW+blockW+pad*3, pad+tabH, editorW-pad, contentH-pad*2, true)
		win32.ShowWindow(a.hEditor, win32.SW_SHOW)
	} else {
		win32.ShowWindow(a.hEditor, win32.SW_HIDE)
	}
}

// -- Terminal Operations --

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
	b.WriteString("--- Command Blocks ---\r\n\r\n")
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
	tcItem := win32.TCITEMW{
		Mask:    win32.TCIF_TEXT,
		PszText: windows.StringToUTF16Ptr(name),
	}
	win32.SendMessage(a.hTabs, win32.TCM_INSERTITEMW, uintptr(a.sessMgr.Count()-1), uintptr(unsafe.Pointer(&tcItem)))

}

func (a *WarpApp) showCommandPalette() {
	var entries []string
	entries = append(entries,
		"/help - Show Warp commands",
		"/ai <query> - Ask AI assistant",
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

// Quit shuts down the application.
func (a *WarpApp) Quit() {
	a.cancel()
	for _, s := range a.sessMgr.List() {
		s.Terminal.Stop()
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
