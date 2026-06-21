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

// AgentStatus represents the state of an agent.
type AgentStatus string

const (
	PendingInit AgentStatus = "PendingInit"
	Running     AgentStatus = "Running"
	Completed   AgentStatus = "Completed"
	Interrupted AgentStatus = "Interrupted"
	Errored     AgentStatus = "Errored"
	Shutdown    AgentStatus = "Shutdown"
)

// TurnContext captures the state for a single turn of execution.
type TurnContext struct {
	TurnID             string
	SessionID          string
	Model              string
	WorkingDir         string
	PermissionsProfile string
}

// ToolCtx holds the execution context.
type ToolCtx struct {
	CallID      string
	ToolName    string
	TurnContext *TurnContext
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

	result := fmt.Sprintf("[%s] Executed shell command '%s' in %s mode (CallID: %s, TurnID: %s)",
		ctx.ToolName, cmdReq, mode, ctx.CallID, ctx.TurnContext.TurnID)
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
		Cwd:         ctx.TurnContext.WorkingDir,
	}

	// 3. Execution
	return runtime.Run(req, attempt, ctx)
}

// AgentSession encapsulates the state machine for processing input.
type AgentSession struct {
	SessionID    string
	Status       AgentStatus
	Orchestrator *ToolOrchestrator
	TurnCounter  int
}

func NewAgentSession(sessionID string) *AgentSession {
	return &AgentSession{
		SessionID:    sessionID,
		Status:       PendingInit,
		Orchestrator: NewToolOrchestrator(),
		TurnCounter:  0,
	}
}

func (a *AgentSession) SteerInput(input string) {
	a.Status = Running
	a.TurnCounter++
	turnID := fmt.Sprintf("turn_%d", a.TurnCounter)

	fmt.Printf("[Agent] Received input: '%s'. Generating TurnContext (%s)...\n", input, turnID)

	turnContext := &TurnContext{
		TurnID:             turnID,
		SessionID:          a.SessionID,
		Model:              "gpt-5.5",
		WorkingDir:         "/workspace",
		PermissionsProfile: "default",
	}

	// Map natural language to a shell tool call for simulation
	toolReq := fmt.Sprintf("echo '%s'", input)
	ctx := &ToolCtx{
		CallID:      fmt.Sprintf("call_%d", a.TurnCounter),
		ToolName:    "shell",
		TurnContext: turnContext,
	}

	runtime := &ShellCommandRuntime{}

	result, err := a.Orchestrator.ExecuteTool(runtime, toolReq, ctx)
	if err != nil {
		fmt.Printf("[Agent] Turn failed. Error: %v\n", err)
		a.Status = Errored
	} else {
		fmt.Printf("[Agent] Turn executed. Result: %v\n", result)
		a.Status = Completed
	}
}

func main() {
	fmt.Println("Welcome to Warp CLI (Go Edition) - Inspired by just-every-code")
	fmt.Println("Type '/help' for commands, or 'quit' to close.")

	agent := NewAgentSession("sess_123")
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
			agent.Status = Shutdown
			break
		}

		handleCommand(input, agent)
	}
}

func handleCommand(input string, agent *AgentSession) {
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
			fmt.Println("  /prompt <msg>- Send a natural language prompt to the agent")
			fmt.Println("  quit         - Quit the application")
		case "prompt":
			if args == "" {
				fmt.Println("[Error] /prompt requires a message.")
			} else {
				agent.SteerInput(args)
			}
		default:
			fmt.Printf("Unknown command: /%s\n", cmd)
		}
	} else {
		agent.SteerInput(input)
	}
}
