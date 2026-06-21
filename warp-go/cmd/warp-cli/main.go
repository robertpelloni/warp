package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// SandboxPreference represents the isolation requirement for a tool.
type SandboxPreference string

const (
	Unsandboxed SandboxPreference = "Unsandboxed"
	Sandboxed   SandboxPreference = "Sandboxed"
)

// ToolCtx holds the execution context.
type ToolCtx struct {
	CallID    string
	ToolName  string
	SessionID string
	TurnID    string
}

// SandboxAttempt represents the environment configuration.
type SandboxAttempt struct {
	IsSandboxed bool
	Cwd         string
}

// ToolRuntime defines the interface tools must implement.
type ToolRuntime interface {
	SandboxPreference() SandboxPreference
	EscalateOnFailure() bool
	RequiresApproval(req interface{}) bool
	Run(req interface{}, attempt *SandboxAttempt, ctx *ToolCtx) (interface{}, error)
}

// ShellCommandRuntime implements ToolRuntime for shell commands.
type ShellCommandRuntime struct{}

func (s *ShellCommandRuntime) SandboxPreference() SandboxPreference {
	return Sandboxed
}

func (s *ShellCommandRuntime) EscalateOnFailure() bool {
	return true
}

func (s *ShellCommandRuntime) RequiresApproval(req interface{}) bool {
	return true
}

func (s *ShellCommandRuntime) Run(req interface{}, attempt *SandboxAttempt, ctx *ToolCtx) (interface{}, error) {
	cmdReq, ok := req.(string)
	if !ok {
		return nil, fmt.Errorf("invalid request type")
	}

	mode := "Unsandboxed"
	if attempt.IsSandboxed {
		mode = "Sandboxed"
	}

	result := fmt.Sprintf("[%s] Executed shell command '%s' in %s mode (CallID: %s)", ctx.ToolName, cmdReq, mode, ctx.CallID)
	return result, nil
}

// ToolOrchestrator coordinates approval, sandboxing, and execution.
type ToolOrchestrator struct{}

func NewToolOrchestrator() *ToolOrchestrator {
	return &ToolOrchestrator{}
}

func (o *ToolOrchestrator) ExecuteTool(runtime ToolRuntime, req interface{}, ctx *ToolCtx) (interface{}, error) {
	// 1. Approval Phase
	if runtime.RequiresApproval(req) {
		fmt.Printf("[Orchestrator] Requesting tool approval for '%s'...\n", ctx.ToolName)
		// Assuming approved
	}

	// 2. Sandbox Selection
	attempt := &SandboxAttempt{
		IsSandboxed: runtime.SandboxPreference() == Sandboxed,
		Cwd:         "/workspace",
	}

	// 3. Execution
	return runtime.Run(req, attempt, ctx)
}

func main() {
	fmt.Println("Welcome to Warp CLI (Go Edition) - Inspired by just-every-code")
	fmt.Println("Type '/help' for commands, or 'quit' to close.")

	orchestrator := NewToolOrchestrator()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("warp> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if strings.EqualFold(input, "quit") || strings.EqualFold(input, "/quit") {
			break
		}

		handleCommand(input, orchestrator)
	}
}

func handleCommand(input string, orchestrator *ToolOrchestrator) {
	if strings.HasPrefix(input, "/") {
		parts := strings.SplitN(input[1:], " ", 2)
		cmd := strings.ToLower(parts[0])
		args := ""
		if len(parts) > 1 {
			args = parts[1]
		}

		switch cmd {
		case "help":
			fmt.Println("Available commands:")
			fmt.Println("  /help        - Show this help message")
			fmt.Println("  /shell <cmd> - Run a command through the ToolOrchestrator")
			fmt.Println("  quit         - Quit the application")
		case "shell":
			if args == "" {
				fmt.Println("[Error] /shell requires a command.")
			} else {
				ctx := &ToolCtx{
					CallID:    "call_abc123",
					ToolName:  "shell",
					SessionID: "sess_1",
					TurnID:    "turn_1",
				}
				runtime := &ShellCommandRuntime{}

				result, err := orchestrator.ExecuteTool(runtime, args, ctx)
				if err != nil {
					fmt.Printf("[Error] %v\n", err)
				} else {
					fmt.Printf("[Result] %v\n", result)
				}
			}
		default:
			fmt.Printf("Unknown command: /%s\n", cmd)
		}
	} else {
		fmt.Printf("[Agent] Echoing input: %s\n", input)
	}
}
