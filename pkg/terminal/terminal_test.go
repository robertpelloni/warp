package terminal

import (
	"strings"
	"testing"
)

func TestParseANSI(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantText   string
		wantColors int
	}{
		{
			name:       "plain text",
			input:      "hello world",
			wantText:   "hello world",
			wantColors: 0,
		},
		{
			name:       "red text",
			input:      "\x1b[31mred text\x1b[0m",
			wantText:   "red text",
			wantColors: 2, // red + reset
		},
		{
			name:       "bold text",
			input:      "\x1b[1mbold\x1b[0m",
			wantText:   "bold",
			wantColors: 2,
		},
		{
			name:       "mixed colors",
			input:      "\x1b[32mgreen\x1b[0m normal \x1b[34mblue\x1b[0m",
			wantText:   "green normal blue",
			wantColors: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, colors := ParseANSI(tt.input)
			if text != tt.wantText {
				t.Errorf("ParseANSI text = %q, want %q", text, tt.wantText)
			}
			if len(colors) != tt.wantColors {
				t.Errorf("ParseANSI colors count = %d, want %d", len(colors), tt.wantColors)
			}
		})
	}
}

func TestNewTerminal(t *testing.T) {
	term := New(Config{
		Shell: "cmd.exe",
		Cols:  80,
		Rows:  24,
	})
	if term == nil {
		t.Fatal("New() returned nil")
	}
	if term.shell != "cmd.exe" {
		t.Errorf("shell = %q, want 'cmd.exe'", term.shell)
	}
}

func TestTerminalSize(t *testing.T) {
	term := New(Config{
		Cols: 120,
		Rows: 40,
	})
	cols, rows := term.Size()
	if cols != 120 || rows != 40 {
		t.Errorf("Size() = (%d, %d), want (120, 40)", cols, rows)
	}
}

func TestTerminalDefaultSize(t *testing.T) {
	term := New(Config{})
	cols, rows := term.Size()
	if cols != 120 || rows != 40 {
		t.Errorf("Default Size() = (%d, %d), want (120, 40)", cols, rows)
	}
}

func TestHistoryNavigation(t *testing.T) {
	term := New(Config{})
	term.history = []string{"ls", "pwd", "cd"}
	term.histIdx = 3

	up := term.HistoryUp()
	if up != "cd" {
		t.Errorf("HistoryUp() = %q, want 'cd'", up)
	}

	up = term.HistoryUp()
	if up != "pwd" {
		t.Errorf("HistoryUp() = %q, want 'pwd'", up)
	}

	down := term.HistoryDown()
	if down != "cd" {
		t.Errorf("HistoryDown() = %q, want 'cd'", down)
	}
}

func TestIsPrompt(t *testing.T) {
	term := New(Config{})

	tests := []struct {
		text     string
		expected bool
	}{
		{"$ ", true},
		{"C:\\> ", true},
		{"# ", true},
		{"hello", false},
		{"", false},
	}

	for _, tt := range tests {
		result := term.isPrompt(tt.text)
		if result != tt.expected {
			t.Errorf("isPrompt(%q) = %v, want %v", tt.text, result, tt.expected)
		}
	}
}

func TestParseANSIColors(t *testing.T) {
	tests := []struct {
		params    string
		final     byte
		wantNil   bool
		wantBold  bool
		wantReset bool
	}{
		{"", 'm', false, false, true},     // empty params with 'm' = reset
		{"1", 'm', false, true, false},    // bold
		{"0", 'm', false, false, true},    // reset
		{"31", 'm', false, false, false},  // red fg
		{"1;31", 'm', false, true, false}, // bold red
		{"x", 'A', true, false, false},    // not 'm' final
	}

	for _, tt := range tests {
		result := parseANSIColor(tt.params, tt.final)
		if tt.wantNil {
			if result != nil {
				t.Errorf("parseANSIColor(%q, %c) expected nil", tt.params, tt.final)
			}
			continue
		}
		if result == nil {
			t.Fatalf("parseANSIColor(%q, %c) returned nil, expected non-nil", tt.params, tt.final)
		}
		if result.Bold != tt.wantBold {
			t.Errorf("parseANSIColor(%q, %c).Bold = %v, want %v", tt.params, tt.final, result.Bold, tt.wantBold)
		}
		if result.Reset != tt.wantReset {
			t.Errorf("parseANSIColor(%q, %c).Reset = %v, want %v", tt.params, tt.final, result.Reset, tt.wantReset)
		}
	}
}

func TestGetOutput(t *testing.T) {
	term := New(Config{})
	if term.GetOutput() != "" {
		t.Error("GetOutput() should be empty for new terminal")
	}
}

func TestHistoryEmpty(t *testing.T) {
	term := New(Config{})
	history := term.History()
	if len(history) != 0 {
		t.Errorf("History() = %d items, want 0", len(history))
	}
}

func TestBlocksInitiallyEmpty(t *testing.T) {
	term := New(Config{})
	blocks := term.Blocks()
	if len(blocks) != 0 {
		t.Errorf("Blocks() = %d items, want 0", len(blocks))
	}
}

// Ensure strings import is used
var _ strings.Builder
