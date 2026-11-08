package terminal

import (
	"bytes"
	"context"
	"fmt"
	"image/color"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/robertpelloni/warp/pkg/pty"
)

// BlockType represents the type of a command block (matching Warp's block model).
type BlockType int

const (
	BlockCommand BlockType = iota // A user-typed command
	BlockOutput                   // Output from a command
	BlockAI                       // AI/agent response
	BlockError                    // Error output
)

// Block represents a Warp-style command block.
type Block struct {
	ID        string
	Type      BlockType
	Command   string
	Output    strings.Builder
	StartTime time.Time
	EndTime   time.Time
	ExitCode  int
	Duration  time.Duration
	Done      bool
}

// Terminal manages the PTY, ANSI parsing, and command block detection.
type Terminal struct {
	mu      sync.RWMutex
	pty     *pty.PTY
	running bool
	cols    int
	rows    int

	// Command block tracking
	blocks       []*Block
	activeBlock  *Block
	promptDetect string

	// ANSI parsing state
	ansiBuf   bytes.Buffer
	escapeBuf bytes.Buffer
	inEscape  bool

	// Callbacks
	OnUpdate   func()
	OnNewBlock func(*Block)

	// Output collector
	outputBuf strings.Builder
	lineBuf   strings.Builder

	// Color state for rendering
	currentFG color.Color
	currentBG color.Color
	bold      bool
	underline bool
	italic    bool

	// Shell detection
	shell string

	// History
	history []string
	histIdx int
}

// Config for creating a Terminal.
type Config struct {
	Shell    string
	Cols     int
	Rows     int
	Dir      string
	OnUpdate func()
}

// New creates a new terminal instance.
func New(cfg Config) *Terminal {
	if cfg.Cols == 0 {
		cfg.Cols = 120
	}
	if cfg.Rows == 0 {
		cfg.Rows = 40
	}
	t := &Terminal{
		cols:         cfg.Cols,
		rows:         cfg.Rows,
		shell:        cfg.Shell,
		promptDetect: "$", // default prompt marker
		OnUpdate:     cfg.OnUpdate,
		history:      make([]string, 0),
	}
	return t
}

// Start launches the PTY and begins processing output.
func (t *Terminal) Start(ctx context.Context) error {
	ptyCfg := pty.Config{
		Shell: t.shell,
		Cols:  uint16(t.cols),
		Rows:  uint16(t.rows),
		Dir:   "",
	}

	p, err := pty.New(ptyCfg)
	if err != nil {
		return fmt.Errorf("create pty: %w", err)
	}
	t.pty = p
	t.running = true

	// Start reading PTY output in background
	go t.readLoop(ctx)

	return nil
}

// Stop closes the PTY.
func (t *Terminal) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.running = false
	if t.pty != nil {
		t.pty.Close()
	}
}

// Write sends input to the PTY.
func (t *Terminal) Write(data []byte) error {
	if t.pty == nil {
		return fmt.Errorf("pty not initialized")
	}
	_, err := t.pty.Write(data)
	return err
}

// WriteString sends a string to the PTY.
func (t *Terminal) WriteString(s string) error {
	return t.Write([]byte(s))
}

// SendCommand sends a command line to the PTY with a trailing newline.
func (t *Terminal) SendCommand(cmd string) error {
	// Create a new command block
	block := &Block{
		ID:        fmt.Sprintf("blk-%d", time.Now().UnixNano()),
		Type:      BlockCommand,
		Command:   cmd,
		StartTime: time.Now(),
	}
	t.mu.Lock()
	if t.activeBlock != nil && !t.activeBlock.Done {
		t.activeBlock.Done = true
		t.activeBlock.EndTime = time.Now()
	}
	t.blocks = append(t.blocks, block)
	t.activeBlock = block
	t.mu.Unlock()

	// Add to history
	t.history = append(t.history, cmd)
	t.histIdx = len(t.history)

	// Notify
	if t.OnNewBlock != nil {
		t.OnNewBlock(block)
	}

	// Send to PTY
	return t.WriteString(cmd + "\n")
}

// Resize changes the terminal dimensions.
func (t *Terminal) Resize(cols, rows int) error {
	t.mu.Lock()
	t.cols = cols
	t.rows = rows
	t.mu.Unlock()
	if t.pty != nil {
		return t.pty.Resize(uint16(cols), uint16(rows))
	}
	return nil
}

// Blocks returns all command blocks.
func (t *Terminal) Blocks() []*Block {
	t.mu.RLock()
	defer t.mu.RUnlock()
	result := make([]*Block, len(t.blocks))
	copy(result, t.blocks)
	return result
}

// ActiveBlock returns the currently executing block.
func (t *Terminal) ActiveBlock() *Block {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.activeBlock
}

// GetOutput returns accumulated terminal output.
func (t *Terminal) GetOutput() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.outputBuf.String()
}

// History returns the command history.
func (t *Terminal) History() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	result := make([]string, len(t.history))
	copy(result, t.history)
	return result
}

// HistoryUp moves up in history and returns the entry.
func (t *Terminal) HistoryUp() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.histIdx > 0 {
		t.histIdx--
	}
	if t.histIdx < len(t.history) {
		return t.history[t.histIdx]
	}
	return ""
}

// HistoryDown moves down in history and returns the entry.
func (t *Terminal) HistoryDown() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.histIdx < len(t.history)-1 {
		t.histIdx++
		return t.history[t.histIdx]
	}
	t.histIdx = len(t.history)
	return ""
}

// Running returns whether the terminal is active.
func (t *Terminal) Running() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.running
}

// Size returns the current terminal dimensions.
func (t *Terminal) Size() (int, int) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.cols, t.rows
}

// Shell returns the shell being used.
func (t *Terminal) Shell() string {
	return t.shell
}

// readLoop continuously reads from the PTY and processes output.
func (t *Terminal) readLoop(ctx context.Context) {
	buf := make([]byte, 4096)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		n, err := t.pty.Read(buf)
		if err != nil {
			if err == io.EOF {
				t.mu.Lock()
				t.running = false
				t.mu.Unlock()
				if t.OnUpdate != nil {
					t.OnUpdate()
				}
				return
			}
			// Non-blocking read returned no data; just wait and retry
			time.Sleep(50 * time.Millisecond)
			continue
		}

		if n > 0 {
			t.processOutput(buf[:n])
		}
	}
}

// processOutput handles raw PTY output, parsing ANSI sequences.
func (t *Terminal) processOutput(data []byte) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.outputBuf.Write(data)

	// Add output to active block
	if t.activeBlock != nil && !t.activeBlock.Done {
		// Try to detect command completion
		text := string(data)
		t.activeBlock.Output.Write(data)

		// Detect prompt to mark block as done
		if t.isPrompt(text) {
			t.activeBlock.Done = true
			t.activeBlock.EndTime = time.Now()
			t.activeBlock.Duration = t.activeBlock.EndTime.Sub(t.activeBlock.StartTime)
		}
	}

	if t.OnUpdate != nil {
		t.OnUpdate()
	}
}

// isPrompt does basic prompt detection.
func (t *Terminal) isPrompt(text string) bool {
	// Heuristic: detect common shell prompts
	// $ or > or # at end of line after a newline
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return false
	}
	last := strings.TrimSpace(lines[len(lines)-1])
	if len(last) == 0 {
		return false
	}
	lastChar := last[len(last)-1]
	return lastChar == '$' || lastChar == '>' || lastChar == '#'
}

// ParseANSI extracts ANSI color information from escape sequences.
// Returns cleaned text and any color changes.
func ParseANSI(data string) (text string, colors []ANSIColor) {
	var buf bytes.Buffer
	var escape bytes.Buffer
	inEscape := false

	for i := 0; i < len(data); i++ {
		ch := data[i]
		if ch == 0x1b && i+1 < len(data) && data[i+1] == '[' {
			inEscape = true
			escape.Reset()
			i++ // skip the '['
			continue
		}
		if inEscape {
			if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
				// End of escape sequence
				seq := escape.String()
				c := parseANSIColor(seq, ch)
				if c != nil {
					colors = append(colors, *c)
				}
				inEscape = false
			} else {
				escape.WriteByte(ch)
			}
		} else {
			buf.WriteByte(ch)
		}
	}

	return buf.String(), colors
}

// ANSIColor represents a color change from an ANSI escape sequence.
type ANSIColor struct {
	FG    color.Color
	BG    color.Color
	Bold  bool
	Reset bool
}

func parseANSIColor(params string, final byte) *ANSIColor {
	if final != 'm' {
		return nil
	}
	c := &ANSIColor{}
	parts := strings.Split(params, ";")
	if len(parts) == 0 || (len(parts) == 1 && parts[0] == "") {
		c.Reset = true
		return c
	}

	for i := 0; i < len(parts); i++ {
		p := parts[i]
		switch p {
		case "0", "":
			c.Reset = true
		case "1":
			c.Bold = true
		case "30":
			c.FG = color.NRGBA{0, 0, 0, 255}
		case "31":
			c.FG = color.NRGBA{205, 0, 0, 255}
		case "32":
			c.FG = color.NRGBA{0, 205, 0, 255}
		case "33":
			c.FG = color.NRGBA{205, 205, 0, 255}
		case "34":
			c.FG = color.NRGBA{0, 0, 238, 255}
		case "35":
			c.FG = color.NRGBA{205, 0, 205, 255}
		case "36":
			c.FG = color.NRGBA{0, 205, 205, 255}
		case "37":
			c.FG = color.NRGBA{229, 229, 229, 255}
		case "38":
			// Extended foreground: 38;5;n (256-color) or 38;2;r;g;b (truecolor)
			if i+2 < len(parts) && parts[i+1] == "5" {
				val, _ := strconv.Atoi(parts[i+2])
				c.FG = parse256Color(val)
				i += 2
			} else if i+4 < len(parts) && parts[i+1] == "2" {
				r, _ := strconv.Atoi(parts[i+2])
				g, _ := strconv.Atoi(parts[i+3])
				b, _ := strconv.Atoi(parts[i+4])
				c.FG = color.NRGBA{uint8(r), uint8(g), uint8(b), 255}
				i += 4
			}
		case "39":
			c.FG = nil // default FG
		case "40":
			c.BG = color.NRGBA{0, 0, 0, 255}
		case "41":
			c.BG = color.NRGBA{205, 0, 0, 255}
		case "42":
			c.BG = color.NRGBA{0, 205, 0, 255}
		case "43":
			c.BG = color.NRGBA{205, 205, 0, 255}
		case "44":
			c.BG = color.NRGBA{0, 0, 238, 255}
		case "45":
			c.BG = color.NRGBA{205, 0, 205, 255}
		case "46":
			c.BG = color.NRGBA{0, 205, 205, 255}
		case "47":
			c.BG = color.NRGBA{229, 229, 229, 255}
		case "48":
			// Extended background: 48;5;n or 48;2;r;g;b
			if i+2 < len(parts) && parts[i+1] == "5" {
				val, _ := strconv.Atoi(parts[i+2])
				c.BG = parse256Color(val)
				i += 2
			} else if i+4 < len(parts) && parts[i+1] == "2" {
				r, _ := strconv.Atoi(parts[i+2])
				g, _ := strconv.Atoi(parts[i+3])
				b, _ := strconv.Atoi(parts[i+4])
				c.BG = color.NRGBA{uint8(r), uint8(g), uint8(b), 255}
				i += 4
			}
		case "49":
			c.BG = nil
		}
	}
	return c
}

func parse256Color(n int) color.Color {
	if n < 0 || n > 255 {
		return color.NRGBA{128, 128, 128, 255}
	}
	// Standard 16 colors
	if n < 16 {
		basic := []color.NRGBA{
			{0, 0, 0, 255}, {205, 0, 0, 255}, {0, 205, 0, 255}, {205, 205, 0, 255},
			{0, 0, 238, 255}, {205, 0, 205, 255}, {0, 205, 205, 255}, {229, 229, 229, 255},
			{127, 127, 127, 255}, {255, 0, 0, 255}, {0, 255, 0, 255}, {255, 255, 0, 255},
			{92, 92, 255, 255}, {255, 0, 255, 255}, {0, 255, 255, 255}, {255, 255, 255, 255},
		}
		return basic[n]
	}
	// 6x6x6 color cube (16-231)
	if n < 232 {
		n -= 16
		r := uint8((n / 36) * 51)
		g := uint8(((n % 36) / 6) * 51)
		b := uint8((n % 6) * 51)
		return color.NRGBA{r, g, b, 255}
	}
	// Grayscale (232-255)
	gray := uint8((n-232)*10 + 8)
	return color.NRGBA{gray, gray, gray, 255}
}

// ─── Virtual Screen Buffer (Milestone 3) ─────────────────────────

// Cell represents a character cell in the virtual terminal grid.
type Cell struct {
	Rune rune
	FG   color.Color
	BG   color.Color
	Bold bool
}

// ScreenBuffer represents a 2D virtual terminal screen buffer.
type ScreenBuffer struct {
	mu      sync.RWMutex
	Cols    int
	Rows    int
	CursorX int
	CursorY int
	Grid    [][]Cell
}

// NewScreenBuffer creates a 2D virtual screen buffer of given dimensions.
func NewScreenBuffer(cols, rows int) *ScreenBuffer {
	if cols <= 0 {
		cols = 80
	}
	if rows <= 0 {
		rows = 24
	}
	sb := &ScreenBuffer{
		Cols: cols,
		Rows: rows,
		Grid: make([][]Cell, rows),
	}
	for r := 0; r < rows; r++ {
		sb.Grid[r] = make([]Cell, cols)
		for c := 0; c < cols; c++ {
			sb.Grid[r][c] = Cell{Rune: ' '}
		}
	}
	return sb
}

// WriteRune writes a rune at current cursor position and advances.
func (sb *ScreenBuffer) WriteRune(r rune, fg, bg color.Color, bold bool) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	if r == '\n' {
		sb.CursorY++
		sb.CursorX = 0
		if sb.CursorY >= sb.Rows {
			sb.scrollUp()
		}
		return
	}
	if r == '\r' {
		sb.CursorX = 0
		return
	}

	if sb.CursorX >= sb.Cols {
		sb.CursorX = 0
		sb.CursorY++
		if sb.CursorY >= sb.Rows {
			sb.scrollUp()
		}
	}

	if sb.CursorY < sb.Rows && sb.CursorX < sb.Cols {
		sb.Grid[sb.CursorY][sb.CursorX] = Cell{
			Rune: r,
			FG:   fg,
			BG:   bg,
			Bold: bold,
		}
		sb.CursorX++
	}
}

// SetCursor moves cursor to 0-based coordinate (col, row).
func (sb *ScreenBuffer) SetCursor(col, row int) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	if col < 0 {
		col = 0
	}
	if col >= sb.Cols {
		col = sb.Cols - 1
	}
	if row < 0 {
		row = 0
	}
	if row >= sb.Rows {
		row = sb.Rows - 1
	}
	sb.CursorX = col
	sb.CursorY = row
}

// MoveCursor relative to current position.
func (sb *ScreenBuffer) MoveCursor(dx, dy int) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	nx := sb.CursorX + dx
	ny := sb.CursorY + dy
	if nx < 0 {
		nx = 0
	}
	if nx >= sb.Cols {
		nx = sb.Cols - 1
	}
	if ny < 0 {
		ny = 0
	}
	if ny >= sb.Rows {
		ny = sb.Rows - 1
	}
	sb.CursorX = nx
	sb.CursorY = ny
}

// Clear clears the entire buffer and resets cursor to (0,0).
func (sb *ScreenBuffer) Clear() {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	for r := 0; r < sb.Rows; r++ {
		for c := 0; c < sb.Cols; c++ {
			sb.Grid[r][c] = Cell{Rune: ' '}
		}
	}
	sb.CursorX = 0
	sb.CursorY = 0
}

func (sb *ScreenBuffer) scrollUp() {
	for r := 1; r < sb.Rows; r++ {
		sb.Grid[r-1] = sb.Grid[r]
	}
	sb.Grid[sb.Rows-1] = make([]Cell, sb.Cols)
	for c := 0; c < sb.Cols; c++ {
		sb.Grid[sb.Rows-1][c] = Cell{Rune: ' '}
	}
	sb.CursorY = sb.Rows - 1
}

// String renders the buffer to a plain text string.
func (sb *ScreenBuffer) String() string {
	sb.mu.RLock()
	defer sb.mu.RUnlock()

	var b strings.Builder
	for r := 0; r < sb.Rows; r++ {
		for c := 0; c < sb.Cols; c++ {
			b.WriteRune(sb.Grid[r][c].Rune)
		}
		if r < sb.Rows-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// ─── Alternate Screen Buffer Manager (for vim, htop, less) ───────

// ScreenBufferManager manages primary and alternate 2D screen buffers.
type ScreenBufferManager struct {
	mu        sync.RWMutex
	Cols      int
	Rows      int
	Primary   *ScreenBuffer
	Alternate *ScreenBuffer
	IsAlt     bool
}

// NewScreenBufferManager creates a screen buffer manager.
func NewScreenBufferManager(cols, rows int) *ScreenBufferManager {
	return &ScreenBufferManager{
		Cols:      cols,
		Rows:      rows,
		Primary:   NewScreenBuffer(cols, rows),
		Alternate: NewScreenBuffer(cols, rows),
		IsAlt:     false,
	}
}

// Active returns the currently displayed screen buffer.
func (m *ScreenBufferManager) Active() *ScreenBuffer {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.IsAlt {
		return m.Alternate
	}
	return m.Primary
}

// EnterAltScreen switches to the alternate screen buffer and blanks it.
func (m *ScreenBufferManager) EnterAltScreen() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.IsAlt = true
	m.Alternate.Clear()
}

// ExitAltScreen switches back to the primary screen buffer.
func (m *ScreenBufferManager) ExitAltScreen() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.IsAlt = false
}

// HandleEscape handles terminal private mode escapes like alternate screen switching.
func (m *ScreenBufferManager) HandleEscape(seq string) bool {
	switch seq {
	case "?1049h", "?47h":
		m.EnterAltScreen()
		return true
	case "?1049l", "?47l":
		m.ExitAltScreen()
		return true
	default:
		return false
	}
}
