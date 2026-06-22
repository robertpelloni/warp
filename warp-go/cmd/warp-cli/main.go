package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

// --- MCP Protocol Definitions ---

type JsonRpcRequest struct {
	ID     string
	Method string
	Params string
}

type JsonRpcResponse struct {
	ID     string
	Result string
	Error  string
}

// --- Orchestration and Context ---

type SandboxPreference string

const (
	Unsandboxed SandboxPreference = "Unsandboxed"
	Sandboxed   SandboxPreference = "Sandboxed"
)

type AgentStatus string

const (
	PendingInit AgentStatus = "PendingInit"
	Running     AgentStatus = "Running"
	Completed   AgentStatus = "Completed"
	Interrupted AgentStatus = "Interrupted"
	Errored     AgentStatus = "Errored"
	Shutdown    AgentStatus = "Shutdown"
)

type TurnContext struct {
	TurnID             string
	SessionID          string
	Model              string
	WorkingDir         string
	PermissionsProfile string
}

type ToolCtx struct {
	CallID      string
	ToolName    string
	TurnContext *TurnContext
}

type SandboxAttempt struct {
	IsSandboxed bool
	Cwd         string
}

type ToolRuntime interface {
	SandboxPreference() SandboxPreference
	EscalateOnFailure() bool
	RequiresApproval(req interface{}) bool
	Run(req interface{}, attempt *SandboxAttempt, ctx *ToolCtx) (interface{}, error)
}

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
	cmdReq, _ := req.(string)
	mode := "Unsandboxed"
	if attempt.IsSandboxed {
		mode = "Sandboxed"
	}
	return fmt.Sprintf("[%s] Executed shell command '%s' in %s mode (CallID: %s, TurnID: %s)",
		ctx.ToolName, cmdReq, mode, ctx.CallID, ctx.TurnContext.TurnID), nil
}

type ToolOrchestrator struct {
	mcpTools map[string]string
	mu       sync.Mutex
}

func NewToolOrchestrator() *ToolOrchestrator {
	return &ToolOrchestrator{
		mcpTools: make(map[string]string),
	}
}

func (o *ToolOrchestrator) RegisterMcpTool(name, description string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.mcpTools[name] = description
	fmt.Printf("[Orchestrator] Registered dynamic MCP tool: %s\n", name)
}

func (o *ToolOrchestrator) ExecuteTool(runtime ToolRuntime, req interface{}, ctx *ToolCtx) (interface{}, error) {
	if runtime.RequiresApproval(req) {
		fmt.Printf("[Orchestrator] Requesting tool approval for '%s'...\n", ctx.ToolName)
	}

	attempt := &SandboxAttempt{
		IsSandboxed: runtime.SandboxPreference() == Sandboxed,
		Cwd:         ctx.TurnContext.WorkingDir,
	}

	return runtime.Run(req, attempt, ctx)
}

// --- MCP Server Implementation ---

type MessageProcessor struct {
	orchestrator *ToolOrchestrator
}

func NewMessageProcessor(orchestrator *ToolOrchestrator) *MessageProcessor {
	return &MessageProcessor{orchestrator: orchestrator}
}

func (m *MessageProcessor) ProcessRequest(req JsonRpcRequest) JsonRpcResponse {
	fmt.Printf("[MCP Server] Processing JSON-RPC method: %s\n", req.Method)

	switch req.Method {
	case "initialize":
		m.orchestrator.RegisterMcpTool("mcp_shell", "Execute commands via MCP")
		return JsonRpcResponse{ID: req.ID, Result: "initialized", Error: ""}
	case "tools/call":
		fmt.Printf("[MCP Server] Dispatched to tool execution via orchestrator: %s\n", req.Params)
		return JsonRpcResponse{ID: req.ID, Result: fmt.Sprintf("Executed MCP tool call with args: %s", req.Params), Error: ""}
	default:
		return JsonRpcResponse{ID: req.ID, Result: "", Error: "Method not found"}
	}
}

// --- Agent Implementation ---

type AgentSession struct {
	SessionID    string
	Status       AgentStatus
	orchestrator *ToolOrchestrator
	mcpProcessor *MessageProcessor
	TurnCounter  int
}

func NewAgentSession(sessionID string) *AgentSession {
	orchestrator := NewToolOrchestrator()
	return &AgentSession{
		SessionID:    sessionID,
		Status:       PendingInit,
		orchestrator: orchestrator,
		mcpProcessor: NewMessageProcessor(orchestrator),
		TurnCounter:  0,
	}
}

func (a *AgentSession) InitializeMcp() {
	req := JsonRpcRequest{ID: "0", Method: "initialize", Params: ""}
	a.mcpProcessor.ProcessRequest(req)
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

	if strings.HasPrefix(input, "mcp") {
		req := JsonRpcRequest{ID: "1", Method: "tools/call", Params: input}
		resp := a.mcpProcessor.ProcessRequest(req)
		fmt.Printf("[Agent] MCP Turn executed. Result: %s\n", resp.Result)
		a.Status = Completed
		return
	}

	toolReq := fmt.Sprintf("echo '%s'", input)
	ctx := &ToolCtx{
		CallID:      fmt.Sprintf("call_%d", a.TurnCounter),
		ToolName:    "shell",
		TurnContext: turnContext,
	}

	runtime := &ShellCommandRuntime{}

	result, err := a.orchestrator.ExecuteTool(runtime, toolReq, ctx)
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
	agent.InitializeMcp()

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
