package terminal

import (
	"bytes"
	"context"
	"fmt"
	"image/color"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/robertpelloni/warp/pkg/pty"
)

// BlockType represents the type of a command block (matching Warp's block model).
type BlockType int

const (
	BlockCommand  BlockType = iota // A user-typed command
	BlockOutput                     // Output from a command
	BlockAI                         // AI/agent response
	BlockError                      // Error output
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
	mu       sync.RWMutex
	pty      *pty.PTY
	running  bool
	cols     int
	rows     int

	// Command block tracking
	blocks     []*Block
	activeBlock *Block
	promptDetect string

	// ANSI parsing state
	ansiBuf     bytes.Buffer
	escapeBuf   bytes.Buffer
	inEscape    bool

	// Callbacks
	OnUpdate    func()
	OnNewBlock  func(*Block)

	// Output collector
	outputBuf   strings.Builder
	lineBuf     strings.Builder

	// Color state for rendering
	currentFG   color.Color
	currentBG   color.Color
	bold        bool
	underline   bool
	italic      bool

	// Shell detection
	shell string

	// History
	history     []string
	histIdx     int
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
		cols:      cfg.Cols,
		rows:      cfg.Rows,
		shell:     cfg.Shell,
		promptDetect: "$", // default prompt marker
		OnUpdate:  cfg.OnUpdate,
		history:   make([]string, 0),
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
		switch parts[i] {
		case "0":
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
		case "49":
			c.BG = nil
		}
	}
	return c
}
