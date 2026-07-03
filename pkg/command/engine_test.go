package command

import (
	"strings"
	"testing"
)

func TestClassify(t *testing.T) {
	e := NewEngine()

	tests := []struct {
		cmd      string
		expected CommandType
	}{
		{"ls -la", CmdShell},
		{"/help", CmdWarp},
		{"/ai what is go", CmdAI},
		{"/warp ai explain this", CmdAI},
		{"ll", CmdAlias},
		{"gs", CmdAlias},
	}

	for _, tt := range tests {
		result := e.Classify(tt.cmd)
		if result != tt.expected {
			t.Errorf("Classify(%q) = %v, want %v", tt.cmd, result, tt.expected)
		}
	}
}

func TestExpand(t *testing.T) {
	e := NewEngine()

	tests := []struct {
		cmd      string
		expected string
	}{
		{"ll", "ls -la"},
		{"gs", "git status"},
		{"gd", "git diff"},
		{"ls -la", "ls -la"}, // not an alias, pass through
		{"echo hello", "echo hello"},
	}

	for _, tt := range tests {
		result := e.Expand(tt.cmd)
		if result != tt.expected {
			t.Errorf("Expand(%q) = %q, want %q", tt.cmd, result, tt.expected)
		}
	}
}

func TestExecuteWarpHelp(t *testing.T) {
	e := NewEngine()
	result, cmdType, err := e.Execute("/help")
	if err != nil {
		t.Fatalf("Execute(/help) error: %v", err)
	}
	if cmdType != CmdWarp {
		t.Errorf("Execute(/help) type = %v, want %v", cmdType, CmdWarp)
	}
	if !strings.Contains(result, "Warp Go") {
		t.Errorf("Execute(/help) result doesn't contain 'Warp Go': %s", result)
	}
}

func TestExecuteAI(t *testing.T) {
	e := NewEngine()
	result, cmdType, err := e.Execute("/ai what is go")
	if err != nil {
		t.Fatalf("Execute(/ai) error: %v", err)
	}
	if cmdType != CmdAI {
		t.Errorf("Execute(/ai) type = %v, want %v", cmdType, CmdAI)
	}
	if !strings.Contains(result, "what is go") {
		t.Errorf("Execute(/ai) result doesn't contain query: %s", result)
	}
}

func TestSuggestions(t *testing.T) {
	e := NewEngine()

	suggestions := e.Suggestions("g")
	if len(suggestions) == 0 {
		t.Error("Suggestions('g') returned no results")
	}

	suggestions = e.Suggestions("/")
	if len(suggestions) == 0 {
		t.Error("Suggestions('/') returned no results")
	}

	suggestions = e.Suggestions("")
	if len(suggestions) != 0 {
		t.Error("Suggestions('') should return empty")
	}
}

func TestRegisterAlias(t *testing.T) {
	e := NewEngine()
	e.RegisterAlias("dcup", "docker compose up -d")

	result := e.Expand("dcup")
	if result != "docker compose up -d" {
		t.Errorf("Expand('dcup') = %q, want 'docker compose up -d'", result)
	}
}

func TestAliases(t *testing.T) {
	e := NewEngine()
	aliases := e.Aliases()
	if len(aliases) == 0 {
		t.Error("Aliases() returned empty map")
	}
	if aliases["ll"] != "ls -la" {
		t.Errorf("Aliases()['ll'] = %q, want 'ls -la'", aliases["ll"])
	}
}

func TestHistory(t *testing.T) {
	e := NewEngine()
	e.Execute("ls")
	e.Execute("pwd")

	history := e.History()
	if len(history) != 2 {
		t.Errorf("History() len = %d, want 2", len(history))
	}
}
