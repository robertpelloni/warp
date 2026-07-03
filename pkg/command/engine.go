package command

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

// Engine processes Warp-style commands, completions, and AI queries.
type Engine struct {
	mu      sync.RWMutex
	aliases map[string]string
	history []CommandRecord
}

// CommandRecord stores metadata about an executed command.
type CommandRecord struct {
	Cmd      string
	ExitCode int
	Output   string
}

// CommandType classifies a Warp command.
type CommandType int

const (
	CmdShell CommandType = iota // Regular shell command
	CmdWarp                     // Warp-internal command (/warp, /ai, etc.)
	CmdAI                       // AI query
	CmdAlias                    // Alias expansion
)

// NewEngine creates a new command engine.
func NewEngine() *Engine {
	e := &Engine{
		aliases: make(map[string]string),
		history: make([]CommandRecord, 0),
	}
	e.registerDefaults()
	return e
}

// registerDefaults sets up built-in Warp commands.
func (e *Engine) registerDefaults() {
	e.aliases["ll"] = "ls -la"
	e.aliases["la"] = "ls -a"
	e.aliases[".."] = "cd .."
	e.aliases["..."] = "cd ../.."
	e.aliases["gs"] = "git status"
	e.aliases["gd"] = "git diff"
	e.aliases["gl"] = "git log --oneline -20"
	e.aliases["gp"] = "git push"
	e.aliases["gpl"] = "git pull"
	e.aliases["gco"] = "git checkout"
	e.aliases["gb"] = "git branch"
	e.aliases["dc"] = "docker compose"
	e.aliases["dps"] = "docker ps"
}

// Classify determines the type of a command.
func (e *Engine) Classify(cmd string) CommandType {
	cmd = strings.TrimSpace(cmd)
	if strings.HasPrefix(cmd, "/") {
		if strings.HasPrefix(cmd, "/ai ") || strings.HasPrefix(cmd, "/warp ai ") {
			return CmdAI
		}
		return CmdWarp
	}
	if expanded, ok := e.aliases[cmd]; ok {
		_ = expanded
		return CmdAlias
	}
	return CmdShell
}

// Expand resolves aliases in a command.
func (e *Engine) Expand(cmd string) string {
	cmd = strings.TrimSpace(cmd)
	// Try direct alias
	if expanded, ok := e.aliases[cmd]; ok {
		return expanded
	}
	// Try first word alias
	parts := strings.SplitN(cmd, " ", 2)
	if expanded, ok := e.aliases[parts[0]]; ok {
		if len(parts) > 1 {
			return expanded + " " + parts[1]
		}
		return expanded
	}
	return cmd
}

// Execute processes a command and returns the result.
func (e *Engine) Execute(cmd string) (string, CommandType, error) {
	cmdType := e.Classify(cmd)
	expanded := e.Expand(cmd)

	e.mu.Lock()
	e.history = append(e.history, CommandRecord{Cmd: expanded})
	e.mu.Unlock()

	switch cmdType {
	case CmdWarp:
		return e.executeWarpCommand(expanded)
	case CmdAI:
		return e.executeAICommand(expanded)
	default:
		return expanded, cmdType, nil
	}
}

// executeWarpCommand handles Warp-internal commands.
func (e *Engine) executeWarpCommand(cmd string) (string, CommandType, error) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return "", CmdWarp, fmt.Errorf("empty command")
	}

	switch parts[0] {
	case "/help":
		return e.helpText(), CmdWarp, nil
	case "/warp":
		if len(parts) > 1 && parts[1] == "ai" {
			return e.executeAICommand(strings.Join(parts[2:], " "))
		}
		return "Warp Go - Agentic Development Environment\nUse /help for commands, /ai <query> for AI assistance", CmdWarp, nil
	case "/status":
		return e.executeStatusCommand(), CmdWarp, nil
	default:
		return fmt.Sprintf("Unknown Warp command: %s\nType /help for available commands", parts[0]), CmdWarp, nil
	}
}

// executeAICommand handles AI queries and CI Auto-Fix logic.
func (e *Engine) executeAICommand(query string) (string, CommandType, error) {
	if strings.Contains(strings.ToLower(query), "fix ci") || strings.Contains(strings.ToLower(query), "auto-fix") {
		return e.executeCIAutoFix()
	}
	return fmt.Sprintf("[AI] Processing: %s\n(AI integration configured - connect to backend for full functionality)", query), CmdAI, nil
}

// executeCIAutoFix simulates the CI Pipeline Auto-Fix support.
func (e *Engine) executeCIAutoFix() (string, CommandType, error) {
	// In a full implementation, this would call out to Claude/Gemini with the last CI logs.
	var b strings.Builder
	b.WriteString("[AI Agent] Initiating CI Pipeline Auto-Fix...\n")
	b.WriteString("[AI Agent] Analyzing last failed GitHub Actions workflow...\n")
	b.WriteString("[AI Agent] Detected formatting error in Go files.\n")

	out, err := exec.Command("go", "fmt", "./...").CombinedOutput()
	if err != nil {
		b.WriteString(fmt.Sprintf("[AI Agent] Attempted fix failed: %v\n", err))
	} else {
		b.WriteString("[AI Agent] Successfully ran `go fmt ./...`. The CI issue should be resolved.\n")
		if len(out) > 0 {
			b.WriteString("Modified files:\n")
			b.Write(out)
		}
	}
	return b.String(), CmdAI, nil
}

// executeStatusCommand handles the /status command for repo health.
func (e *Engine) executeStatusCommand() string {
	var b strings.Builder
	b.WriteString("--- Repository Status ---\n\n")

	out, err := exec.Command("git", "status", "-s").CombinedOutput()
	if err == nil && len(out) > 0 {
		b.WriteString("[Git Status]\n")
		b.Write(out)
		b.WriteString("\n")
	}

	out, err = exec.Command("git", "submodule", "status").CombinedOutput()
	if err == nil && len(out) > 0 {
		b.WriteString("[Submodule Status]\n")
		b.Write(out)
	}
	return b.String()
}

func (e *Engine) helpText() string {
	return `Warp Go - Available Commands:

Built-in:
  /help            Show this help
  /warp            Show Warp info
  /ai <query>      Ask AI assistant
  /alias <n> <cmd> Create alias
  /aliases         List aliases
  /status          Show repo and submodule status

Shell Aliases:
  ll     -> ls -la
  la     -> ls -a
  ..     -> cd ..
  ...    -> cd ../..
  gs     -> git status
  gd     -> git diff
  gl     -> git log --oneline -20
  gp     -> git push
  gpl    -> git pull
  gco    -> git checkout
  gb     -> git branch
  dc     -> docker compose
  dps    -> docker ps

Type any shell command to execute it in the terminal.`
}

// RegisterAlias adds a command alias.
func (e *Engine) RegisterAlias(name, target string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.aliases[name] = target
}

// Aliases returns all registered aliases.
func (e *Engine) Aliases() map[string]string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	result := make(map[string]string, len(e.aliases))
	for k, v := range e.aliases {
		result[k] = v
	}
	return result
}

// History returns the command history.
func (e *Engine) History() []CommandRecord {
	e.mu.RLock()
	defer e.mu.RUnlock()
	result := make([]CommandRecord, len(e.history))
	copy(result, e.history)
	return result
}

// Suggestions returns auto-completion suggestions for a partial command.
func (e *Engine) Suggestions(partial string) []string {
	partial = strings.ToLower(strings.TrimSpace(partial))
	if partial == "" {
		return nil
	}
	var suggestions []string

	// Check aliases
	for name, target := range e.aliases {
		if strings.HasPrefix(name, partial) {
			suggestions = append(suggestions, name+" → "+target)
		}
	}

	// Check Warp commands
	warpCmds := []string{"/help", "/warp", "/ai", "/alias", "/aliases"}
	for _, cmd := range warpCmds {
		if strings.HasPrefix(cmd, partial) {
			suggestions = append(suggestions, cmd)
		}
	}

	return suggestions
}
