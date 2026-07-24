//go:build windows

package win32

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	moduser32   = windows.NewLazyDLL("user32.dll")
	modkernel32 = windows.NewLazyDLL("kernel32.dll")
	modgdi32    = windows.NewLazyDLL("gdi32.dll")
	moduxtheme  = windows.NewLazyDLL("uxtheme.dll")
	modcomctl32 = windows.NewLazyDLL("comctl32.dll")

	procCreateWindowExW       = moduser32.NewProc("CreateWindowExW")
	procDefWindowProcW        = moduser32.NewProc("DefWindowProcW")
	procRegisterClassExW      = moduser32.NewProc("RegisterClassExW")
	procGetMessageW           = moduser32.NewProc("GetMessageW")
	procTranslateMessage      = moduser32.NewProc("TranslateMessage")
	procDispatchMessageW      = moduser32.NewProc("DispatchMessageW")
	procPostQuitMessage       = moduser32.NewProc("PostQuitMessage")
	procShowWindow            = moduser32.NewProc("ShowWindow")
	procUpdateWindow          = moduser32.NewProc("UpdateWindow")
	procSetWindowTextW        = moduser32.NewProc("SetWindowTextW")
	procGetWindowTextW        = moduser32.NewProc("GetWindowTextW")
	procGetWindowTextLengthW  = moduser32.NewProc("GetWindowTextLengthW")
	procSendMessageW          = moduser32.NewProc("SendMessageW")
	procPostMessageW          = moduser32.NewProc("PostMessageW")
	procSetWindowPos          = moduser32.NewProc("SetWindowPos")
	procGetClientRect         = moduser32.NewProc("GetClientRect")
	procMoveWindow            = moduser32.NewProc("MoveWindow")
	procInvalidateRect        = moduser32.NewProc("InvalidateRect")
	procSetFocus              = moduser32.NewProc("SetFocus")
	procSetForegroundWindow   = moduser32.NewProc("SetForegroundWindow")
	procDestroyWindow         = moduser32.NewProc("DestroyWindow")
	procGetDC                 = moduser32.NewProc("GetDC")
	procReleaseDC             = moduser32.NewProc("ReleaseDC")
	procBeginPaint            = moduser32.NewProc("BeginPaint")
	procEndPaint              = moduser32.NewProc("EndPaint")
	procFillRect              = moduser32.NewProc("FillRect")
	procDrawTextW             = moduser32.NewProc("DrawTextW")
	procSetTextColor          = modgdi32.NewProc("SetTextColor")
	procSetBkColor            = modgdi32.NewProc("SetBkColor")
	procSetBkMode             = modgdi32.NewProc("SetBkMode")
	procCreateSolidBrush      = modgdi32.NewProc("CreateSolidBrush")
	procCreateFontW           = modgdi32.NewProc("CreateFontW")
	procSelectObject          = modgdi32.NewProc("SelectObject")
	procDeleteObject          = modgdi32.NewProc("DeleteObject")
	procGetStockObject        = modgdi32.NewProc("GetStockObject")
	procSetTimer              = moduser32.NewProc("SetTimer")
	procKillTimer             = moduser32.NewProc("KillTimer")
	procMessageBeep           = moduser32.NewProc("MessageBeep")
	procSetWindowTextA        = moduser32.NewProc("SetWindowTextA")
	procGetSysColor           = moduser32.NewProc("GetSysColor")
	procLoadCursorW           = moduser32.NewProc("LoadCursorW")
	procLoadIconW             = moduser32.NewProc("LoadIconW")
	procInitCommonControlsEx  = modcomctl32.NewProc("InitCommonControlsEx")
	procSetWindowText         = moduser32.NewProc("SetWindowTextW")
	procGetDlgItem            = moduser32.NewProc("GetDlgItem")
	procCreateFontIndirectW   = modgdi32.NewProc("CreateFontIndirectW")
	procTextOutW              = modgdi32.NewProc("TextOutW")
	procGetTextExtentPoint32W = modgdi32.NewProc("GetTextExtentPoint32W")
	procSetDCPenColor         = modgdi32.NewProc("SetDCPenColor")
	procRectangle             = modgdi32.NewProc("Rectangle")
	procLineTo                = modgdi32.NewProc("LineTo")
	procMoveToEx              = modgdi32.NewProc("MoveToEx")
	procCreatePen             = modgdi32.NewProc("CreatePen")
)

// Window handles and types
type HWND uintptr
type HINSTANCE uintptr
type HFONT uintptr
type HBRUSH uintptr
type HDC uintptr
type HCURSOR uintptr
type HICON uintptr

// Window style constants
const (
	WS_OVERLAPPED       = 0x00000000
	WS_POPUP            = 0x80000000
	WS_CHILD            = 0x40000000
	WS_MINIMIZE         = 0x20000000
	WS_VISIBLE          = 0x10000000
	WS_DISABLED         = 0x08000000
	WS_CLIPSIBLINGS     = 0x04000000
	WS_CLIPCHILDREN     = 0x02000000
	WS_MAXIMIZE         = 0x01000000
	WS_CAPTION          = 0x00C00000
	WS_BORDER           = 0x00800000
	WS_DLGFRAME         = 0x00400000
	WS_VSCROLL          = 0x00200000
	WS_HSCROLL          = 0x00100000
	WS_SYSMENU          = 0x00080000
	WS_THICKFRAME       = 0x00040000
	WS_MINIMIZEBOX      = 0x00020000
	WS_MAXIMIZEBOX      = 0x00010000
	WS_OVERLAPPEDWINDOW = WS_OVERLAPPED | WS_CAPTION | WS_SYSMENU | WS_THICKFRAME | WS_MINIMIZEBOX | WS_MAXIMIZEBOX
	WS_EX_APPWINDOW     = 0x00040000
	WS_EX_CLIENTEDGE    = 0x00000200

	// Edit control styles
	ES_MULTILINE     = 0x0004
	ES_AUTOVSCROLL   = 0x0040
	ES_AUTOHSCROLL   = 0x0080
	ES_WANTRETURN    = 0x1000
	ES_READONLY      = 0x0800

	// Show window constants
	SW_SHOW      = 5
	SW_SHOWDEFAULT = 10

	// Messages
	WM_DESTROY       = 0x0002
	WM_PAINT         = 0x000F
	WM_COMMAND       = 0x0111
	WM_SIZE          = 0x0005
	WM_SETFONT       = 0x0030
	WM_TIMER         = 0x0113
	WM_KEYDOWN       = 0x0100
	WM_KEYUP         = 0x0101
	WM_CHAR          = 0x0102
	WM_ERASEBKGND    = 0x0014
	WM_CTLCOLOREDIT  = 0x0133
	WM_CTLCOLORSTATIC = 0x0138
	WM_GETMINMAXINFO  = 0x0024
	WM_NCCREATE       = 0x0081
	WM_CREATE         = 0x0001
	WM_LBUTTONDOWN    = 0x0201
	WM_SETTEXT        = 0x000C
	WM_GETTEXT        = 0x000D
	WM_GETTEXTLENGTH  = 0x000E
	EM_SETSEL        = 0x00B1
	EM_REPLACESEL    = 0x00C2
	EM_GETLINECOUNT  = 0x00BA
	EM_SCROLL        = 0x00B5
	EM_LINESCROLL    = 0x00B6
	EM_SETREADONLY   = 0x00CF

	// Scroll bar constants
	SB_VSCROLL     = 1
	SB_BOTTOM      = 7
	SB_LINEDOWN    = 1
	SB_LINEUP      = 0

	// Menu constants
	MF_STRING      = 0x00000000
	MF_SEPARATOR   = 0x00000800
	MF_POPUP       = 0x00000010

	// Virtual key codes
	VK_RETURN      = 0x0D
	VK_UP          = 0x26
	VK_DOWN        = 0x28
	VK_CONTROL     = 0x11
	VK_SHIFT       = 0x10
	VK_TAB         = 0x09
	VK_ESCAPE      = 0x1B
	VK_BACK        = 0x08
	VK_HOME        = 0x24
	VK_END         = 0x23
	VK_PRIOR       = 0x21 // Page Up
	VK_NEXT        = 0x22 // Page Down
	VK_F5          = 0x74

	// Color constants
	COLOR_WINDOW     = 5
	COLOR_WINDOWTEXT = 8
	COLOR_BTNFACE    = 15
	TRANSPARENT      = 1
	OPAQUE           = 2

	// DrawText flags
	DT_LEFT       = 0x0000
	DT_TOP        = 0x0000
	DT_WORDBREAK  = 0x0010
	DT_EXPANDTABS = 0x0040
	DT_NOPREFIX   = 0x0200

	// Font constants
	DEFAULT_CHARSET     = 1
	OUT_DEFAULT_PRECIS  = 0
	CLIP_DEFAULT_PRECIS = 0
	DEFAULT_QUALITY     = 0
	DEFAULT_PITCH       = 0
	FW_NORMAL           = 400
	FW_BOLD             = 700
	FIXED_PITCH         = 1
	FF_DONTCARE         = 0
	FF_MODERN           = 3

	// Timer
	IDT_REFRESH = 1
)

// MSG is the Windows message structure.
type MSG struct {
	HWnd    HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

// POINT is a Windows POINT structure.
type POINT struct {
	X, Y int32
}

// RECT is a Windows RECT structure.
type RECT struct {
	Left, Top, Right, Bottom int32
}

// SIZE is a Windows SIZE structure.
type SIZE struct {
	CX, CY int32
}

// WNDCLASSEXW is the extended window class structure.
type WNDCLASSEXW struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     HINSTANCE
	HIcon         HICON
	HCursor       HCURSOR
	HBrBackground HBRUSH
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       HICON
}

// PAINTSTRUCT is the paint structure.
type PAINTSTRUCT struct {
	HDC         HDC
	FErase      int32
	RcPaint     RECT
	FRestore    int32
	FIncUpdate  int32
	RgbReserved [32]byte
}

// CREATESTRUCTW is passed with WM_CREATE.
type CREATESTRUCTW struct {
	LpCreateParams uintptr
	HInstance      HINSTANCE
	HMenu          uintptr
	HWndParent     HWND
	CY             int32
	CX             int32
	Y              int32
	X              int32
	Style          int32
	LpszName       *uint16
	LpszClass      *uint16
	ExStyle        uint32
}

// MINMAXINFO is passed with WM_GETMINMAXINFO.
type MINMAXINFO struct {
	PtReserved     POINT
	PtMaxSize      POINT
	PtMaxPosition  POINT
	PtMinTrackSize POINT
	PtMaxTrackSize POINT
}

// LOGFONTW describes a font.
type LOGFONTW struct {
	LfHeight         int32
	LfWidth          int32
	LfEscapement     int32
	LfOrientation    int32
	LfWeight         int32
	LfItalic         byte
	LfUnderline      byte
	LfStrikeOut      byte
	LfCharSet        byte
	LfOutPrecision   byte
	LfClipPrecision  byte
	LfQuality        byte
	LfPitchAndFamily byte
	LfFaceName       [32]uint16
}

// WindowProc is the window procedure callback type.
type WindowProc func(hWnd HWND, msg uint32, wParam uintptr, lParam uintptr) uintptr

// CreateWindowEx creates a window.
func CreateWindowEx(exStyle uint32, className, windowName *uint16, style uint32, x, y, width, height int32, parent HWND, menu uintptr, instance HINSTANCE, param uintptr) HWND {
	ret, _, _ := procCreateWindowExW.Call(
		uintptr(exStyle),
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)),
		uintptr(style),
		uintptr(x), uintptr(y),
		uintptr(width), uintptr(height),
		uintptr(parent),
		menu,
		uintptr(instance),
		param,
	)
	return HWND(ret)
}

// RegisterClassEx registers a window class.
func RegisterClassEx(wc *WNDCLASSEXW) uint16 {
	ret, _, _ := procRegisterClassExW.Call(uintptr(unsafe.Pointer(wc)))
	return uint16(ret)
}

// DefWindowProc calls the default window procedure.
func DefWindowProc(hWnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := procDefWindowProcW.Call(uintptr(hWnd), uintptr(msg), wParam, lParam)
	return ret
}

// GetMessage retrieves a message.
func GetMessage(msg *MSG, hWnd HWND, wMsgFilterMin, wMsgFilterMax uint32) int32 {
	ret, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(msg)), uintptr(hWnd), uintptr(wMsgFilterMin), uintptr(wMsgFilterMax))
	return int32(ret)
}

// TranslateMessage translates virtual-key messages.
func TranslateMessage(msg *MSG) bool {
	ret, _, _ := procTranslateMessage.Call(uintptr(unsafe.Pointer(msg)))
	return ret != 0
}

// DispatchMessage dispatches a message.
func DispatchMessage(msg *MSG) uintptr {
	ret, _, _ := procDispatchMessageW.Call(uintptr(unsafe.Pointer(msg)))
	return ret
}

// PostQuitMessage posts a quit message.
func PostQuitMessage(exitCode int32) {
	procPostQuitMessage.Call(uintptr(exitCode))
}

// ShowWindow shows a window.
func ShowWindow(hWnd HWND, cmdShow int32) bool {
	ret, _, _ := procShowWindow.Call(uintptr(hWnd), uintptr(cmdShow))
	return ret != 0
}

// UpdateWindow updates a window.
func UpdateWindow(hWnd HWND) {
	procUpdateWindow.Call(uintptr(hWnd))
}

// SetWindowText sets window text.
func SetWindowText(hWnd HWND, text string) {
	procSetWindowText.Call(uintptr(hWnd), uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(text))))
}

// GetWindowText gets window text.
func GetWindowText(hWnd HWND) string {
	len := GetWindowTextLength(hWnd)
	if len == 0 {
		return ""
	}
	buf := make([]uint16, len+1)
	procGetWindowTextW.Call(
		uintptr(hWnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len+1),
	)
	return windows.UTF16ToString(buf)
}

// GetWindowTextLength gets the length of window text.
func GetWindowTextLength(hWnd HWND) int32 {
	ret, _, _ := procGetWindowTextLengthW.Call(uintptr(hWnd))
	return int32(ret)
}

// SendMessage sends a message.
func SendMessage(hWnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := procSendMessageW.Call(uintptr(hWnd), uintptr(msg), wParam, lParam)
	return ret
}

// PostMessage posts a message.
func PostMessage(hWnd HWND, msg uint32, wParam, lParam uintptr) bool {
	ret, _, _ := procPostMessageW.Call(uintptr(hWnd), uintptr(msg), wParam, lParam)
	return ret != 0
}

// GetClientRect gets the client area rectangle.
func GetClientRect(hWnd HWND) RECT {
	var rc RECT
	procGetClientRect.Call(uintptr(hWnd), uintptr(unsafe.Pointer(&rc)))
	return rc
}

// MoveWindow repositions and resizes a window.
func MoveWindow(hWnd HWND, x, y, width, height int32, repaint bool) {
	procMoveWindow.Call(uintptr(hWnd), uintptr(x), uintptr(y), uintptr(width), uintptr(height), boolToPtr(repaint))
}

// InvalidateRect invalidates a rectangle.
func InvalidateRect(hWnd HWND, rc *RECT, erase bool) {
	procInvalidateRect.Call(uintptr(hWnd), uintptr(unsafe.Pointer(rc)), boolToPtr(erase))
}

// SetFocus sets keyboard focus.
func SetFocus(hWnd HWND) HWND {
	ret, _, _ := procSetFocus.Call(uintptr(hWnd))
	return HWND(ret)
}

// DestroyWindow destroys a window.
func DestroyWindow(hWnd HWND) bool {
	ret, _, _ := procDestroyWindow.Call(uintptr(hWnd))
	return ret != 0
}

// GetDC gets a device context.
func GetDC(hWnd HWND) HDC {
	ret, _, _ := procGetDC.Call(uintptr(hWnd))
	return HDC(ret)
}

// ReleaseDC releases a device context.
func ReleaseDC(hWnd HWND, hdc HDC) int32 {
	ret, _, _ := procReleaseDC.Call(uintptr(hWnd), uintptr(hdc))
	return int32(ret)
}

// BeginPaint prepares for painting.
func BeginPaint(hWnd HWND, ps *PAINTSTRUCT) HDC {
	ret, _, _ := procBeginPaint.Call(uintptr(hWnd), uintptr(unsafe.Pointer(ps)))
	return HDC(ret)
}

// EndPaint ends painting.
func EndPaint(hWnd HWND, ps *PAINTSTRUCT) {
	procEndPaint.Call(uintptr(hWnd), uintptr(unsafe.Pointer(ps)))
}

// FillRect fills a rectangle.
func FillRect(hdc HDC, rc *RECT, brush HBRUSH) int32 {
	ret, _, _ := procFillRect.Call(uintptr(hdc), uintptr(unsafe.Pointer(rc)), uintptr(brush))
	return int32(ret)
}

// DrawText draws formatted text.
func DrawText(hdc HDC, text string, rc *RECT, format uint32) int32 {
	procDrawTextW.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(text))),
		1<<32 - 1,
		uintptr(unsafe.Pointer(rc)),
		uintptr(format),
	)
	// Return value is height
	return 0
}

// SetTextColor sets the text color.
func SetTextColor(hdc HDC, color uint32) uint32 {
	ret, _, _ := procSetTextColor.Call(uintptr(hdc), uintptr(color))
	return uint32(ret)
}

// SetBkColor sets the background color.
func SetBkColor(hdc HDC, color uint32) uint32 {
	ret, _, _ := procSetBkColor.Call(uintptr(hdc), uintptr(color))
	return uint32(ret)
}

// SetBkMode sets the background mode.
func SetBkMode(hdc HDC, mode int32) int32 {
	ret, _, _ := procSetBkMode.Call(uintptr(hdc), uintptr(mode))
	return int32(ret)
}

// CreateSolidBrush creates a solid color brush.
func CreateSolidBrush(color uint32) HBRUSH {
	ret, _, _ := procCreateSolidBrush.Call(uintptr(color))
	return HBRUSH(ret)
}

// CreateFont creates a font.
func CreateFont(height, width int32, escapement, orientation, weight int32, italic, underline, strikeOut, charSet, outPrecision, clipPrecision, quality, pitchAndFamily byte, faceName string) HFONT {
	ret, _, _ := procCreateFontW.Call(
		uintptr(height), uintptr(width),
		uintptr(escapement), uintptr(orientation),
		uintptr(weight),
		uintptr(italic), uintptr(underline), uintptr(strikeOut),
		uintptr(charSet), uintptr(outPrecision), uintptr(clipPrecision),
		uintptr(quality), uintptr(pitchAndFamily),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(faceName))),
	)
	return HFONT(ret)
}

// SelectObject selects a GDI object.
func SelectObject(hdc HDC, obj uintptr) uintptr {
	ret, _, _ := procSelectObject.Call(uintptr(hdc), obj)
	return ret
}

// DeleteObject deletes a GDI object.
func DeleteObject(obj uintptr) bool {
	ret, _, _ := procDeleteObject.Call(obj)
	return ret != 0
}

// SetTimer sets a timer.
func SetTimer(hWnd HWND, nIDEvent uintptr, elapse uint32, timerProc uintptr) uintptr {
	ret, _, _ := procSetTimer.Call(uintptr(hWnd), nIDEvent, uintptr(elapse), timerProc)
	return ret
}

// KillTimer destroys a timer.
func KillTimer(hWnd HWND, nIDEvent uintptr) bool {
	ret, _, _ := procKillTimer.Call(uintptr(hWnd), nIDEvent)
	return ret != 0
}

// LoadCursor loads a cursor.
func LoadCursor(hInstance HINSTANCE, cursor uintptr) HCURSOR {
	ret, _, _ := procLoadCursorW.Call(uintptr(hInstance), cursor)
	return HCURSOR(ret)
}

// LoadIcon loads an icon.
func LoadIcon(hInstance HINSTANCE, icon uintptr) HICON {
	ret, _, _ := procLoadIconW.Call(uintptr(hInstance), icon)
	return HICON(ret)
}

// RGB creates a COLORREF from R, G, B values.
func RGB(r, g, b uint8) uint32 {
	return uint32(r) | uint32(g)<<8 | uint32(b)<<16
}

// NRGBToCOLORREF converts an NRGBA color to a COLORREF.
func NRGBToCOLORREF(r, g, b, a uint8) uint32 {
	_ = a // COLORREF doesn't have alpha
	return RGB(r, g, b)
}

// GetModuleHandle gets the current module handle.
func GetModuleHandle() HINSTANCE {
	ret, _, _ := modkernel32.NewProc("GetModuleHandleW").Call(0)
	return HINSTANCE(ret)
}

// IsKeyDown checks if a virtual key is down.
func IsKeyDown(vk int) bool {
	ret, _, _ := moduser32.NewProc("GetAsyncKeyState").Call(uintptr(vk))
	return ret&0x8000 != 0
}

// HiWord returns the high word of a uintptr.
func HiWord(lParam uintptr) uint16 {
	return uint16(lParam >> 16)
}

// LoWord returns the low word of a uintptr.
func LoWord(lParam uintptr) uint16 {
	return uint16(lParam & 0xFFFF)
}

// MakeLParam makes an lParam from x, y coordinates.
func MakeLParam(x, y int16) uintptr {
	return uintptr(uint16(x)) | uintptr(uint16(y))<<16
}

// MakeWParam makes a wParam from lo, hi values.
func MakeWParam(lo, hi uint16) uintptr {
	return uintptr(lo) | uintptr(hi)<<16
}

func boolToPtr(b bool) uintptr {
	if b {
		return 1
	}
	return 0
}

// SetWindowLongPtr sets a window long value.
func SetWindowLongPtr(hWnd HWND, nIndex int32, dwNewLong uintptr) uintptr {
	ret, _, _ := moduser32.NewProc("SetWindowLongPtrW").Call(uintptr(hWnd), uintptr(nIndex), dwNewLong)
	return ret
}

// GetWindowLongPtr gets a window long value.
func GetWindowLongPtr(hWnd HWND, nIndex int32) uintptr {
	ret, _, _ := moduser32.NewProc("GetWindowLongPtrW").Call(uintptr(hWnd), uintptr(nIndex))
	return ret
}

// CallWindowProc calls a window procedure.
func CallWindowProc(lpPrevWndFunc uintptr, hWnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := moduser32.NewProc("CallWindowProcW").Call(lpPrevWndFunc, uintptr(hWnd), uintptr(msg), wParam, lParam)
	return ret
}

// GWLP_WNDPROC for SetWindowLongPtr.
const GWLP_WNDPROC = -4

// Debug placeholder
var _ = fmt.Sprintf
