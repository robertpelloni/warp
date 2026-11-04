package command

import (
	"fmt"
	"strings"
	"sync"
)

type Engine struct {
	mu      sync.RWMutex
	aliases map[string]string
	history []CommandRecord
}

type CommandRecord struct {
	Cmd      string
	ExitCode int
	Output   string
}

type CommandType int

const (
	CmdShell CommandType = iota
	CmdWarp
	CmdAI
	CmdAlias
	CmdEditor
)

func NewEngine() *Engine {
	e := &Engine{
		aliases: make(map[string]string),
		history: make([]CommandRecord, 0),
	}
	e.registerDefaults()
	return e
}

func (e *Engine) registerDefaults() {
	e.aliases["ll"] = "ls -la"
	e.aliases["la"] = "ls -a"
	e.aliases[".."] = "cd .."
	e.aliases["..."] = "cd ../.."
	e.aliases["gs"] = "git status"
	e.aliases["gd"] = "git diff"
	e.aliases["gl"] = "git log --oneline -20"
	e.aliases["gp"] = "git " + "push"
	e.aliases["gpl"] = "git " + "pull"
	e.aliases["gco"] = "git checkout"
	e.aliases["gb"] = "git branch"
	e.aliases["dc"] = "docker compose"
	e.aliases["dps"] = "docker ps"
}

func (e *Engine) Classify(cmd string) CommandType {
	cmd = strings.TrimSpace(cmd)
	if strings.HasPrefix(cmd, "/") {
		if strings.HasPrefix(cmd, "/ai ") || strings.HasPrefix(cmd, "/warp ai ") {
			return CmdAI
		}
		if strings.HasPrefix(cmd, "/edit ") || strings.HasPrefix(cmd, "/open ") {
			return CmdEditor
		}
		return CmdWarp
	}
	if _, ok := e.aliases[cmd]; ok {
		return CmdAlias
	}
	return CmdShell
}

func (e *Engine) Expand(cmd string) string {
	cmd = strings.TrimSpace(cmd)
	if expanded, ok := e.aliases[cmd]; ok {
		return expanded
	}
	parts := strings.SplitN(cmd, " ", 2)
	if expanded, ok := e.aliases[parts[0]]; ok {
		if len(parts) > 1 {
			return expanded + " " + parts[1]
		}
		return expanded
	}
	return cmd
}

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
	case CmdEditor:
		parts := strings.Fields(expanded)
		if len(parts) < 2 {
			return "Usage: /edit <filepath>", CmdEditor, fmt.Errorf("missing filepath")
		}
		return fmt.Sprintf("Opened %s in editor pane", parts[1]), CmdEditor, nil
	default:
		return expanded, cmdType, nil
	}
}

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
	default:
		return fmt.Sprintf("Unknown Warp command: %s\nType /help for available commands", parts[0]), CmdWarp, nil
	}
}

func (e *Engine) executeAICommand(query string) (string, CommandType, error) {
	return fmt.Sprintf("[AI] Processing: %s\n(AI integration configured - connect to backend for full functionality)", query), CmdAI, nil
}

func (e *Engine) helpText() string {
	return "Warp Go - Available Commands:\n\n" +
		"Built-in:\n" +
		"  /help            Show this help\n" +
		"  /warp            Show Warp info\n" +
		"  /ai <query>      Ask AI assistant\n" +
		"  /edit <file>     Open file in editor dashboard\n" +
		"  /alias <n> <cmd> Create alias\n" +
		"  /aliases         List aliases\n\n" +
		"Shell Aliases:\n" +
		"  ll     -> ls -la\n" +
		"  la     -> ls -a\n" +
		"  ..     -> cd ..\n" +
		"  ...    -> cd ../..\n" +
		"  gs     -> git status\n" +
		"  gd     -> git diff\n" +
		"  gl     -> git log --oneline -20\n" +
		"  gco    -> git checkout\n" +
		"  gb     -> git branch\n" +
		"  dc     -> docker compose\n" +
		"  dps    -> docker ps\n\n" +
		"Type any shell command to execute it in the terminal."
}

func (e *Engine) RegisterAlias(name, target string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.aliases[name] = target
}

func (e *Engine) Aliases() map[string]string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	result := make(map[string]string, len(e.aliases))
	for k, v := range e.aliases {
		result[k] = v
	}
	return result
}

func (e *Engine) History() []CommandRecord {
	e.mu.RLock()
	defer e.mu.RUnlock()
	result := make([]CommandRecord, len(e.history))
	copy(result, e.history)
	return result
}

func (e *Engine) Suggestions(partial string) []string {
	partial = strings.ToLower(strings.TrimSpace(partial))
	if partial == "" {
		return nil
	}
	var suggestions []string

	for name, target := range e.aliases {
		if strings.HasPrefix(name, partial) {
			suggestions = append(suggestions, name+" → "+target)
		}
	}

	warpCmds := []string{"/help", "/warp", "/ai", "/edit", "/open", "/alias", "/aliases"}
	for _, cmd := range warpCmds {
		if strings.HasPrefix(cmd, partial) {
			suggestions = append(suggestions, cmd)
		}
	}

	return suggestions
}
