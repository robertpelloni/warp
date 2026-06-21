package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type ToolOrchestrator struct {
	SandboxEnabled bool
}

func NewToolOrchestrator() *ToolOrchestrator {
	return &ToolOrchestrator{SandboxEnabled: true}
}

func (o *ToolOrchestrator) ExecuteTask(task string) string {
	fmt.Printf("[Orchestrator] Requesting approval for task: '%s'\n", task)
	time.Sleep(500 * time.Millisecond)

	if o.SandboxEnabled {
		fmt.Println("[Orchestrator] Running in sandbox mode...")
	}

	steps := []string{"Plan", "Code", "Review"}
	for _, step := range steps {
		fmt.Printf("[Agent: %s] Processing...\n", step)
		time.Sleep(400 * time.Millisecond)
	}

	return fmt.Sprintf("Task '%s' completed successfully by Auto Drive.", task)
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
			fmt.Println("  /auto <task> - Start Auto Drive orchestration for a task")
			fmt.Println("  quit         - Quit the application")
		case "auto":
			if args == "" {
				fmt.Println("[Error] /auto requires a task description.")
			} else {
				fmt.Println("[Orchestrator: Auto Drive] Starting autonomous loop.")
				result := orchestrator.ExecuteTask(args)
				fmt.Printf("[Result] %s\n", result)
			}
		default:
			fmt.Printf("Unknown command: /%s\n", cmd)
		}
	} else {
		fmt.Printf("[Agent] Echoing input: %s\n", input)
	}
}
