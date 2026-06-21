package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Welcome to Warp CLI (Go Edition) - Inspired by just-every-code")
	fmt.Println("Type '/help' for commands, or 'quit' to close.")

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

		if input == "quit" || input == "/quit" {
			break
		}

		handleCommand(input)
	}
}

func handleCommand(input string) {
	if strings.HasPrefix(input, "/") {
		cmd := strings.TrimPrefix(input, "/")
		switch cmd {
		case "help":
			fmt.Println("Available commands:")
			fmt.Println("  /help     - Show this help message")
			fmt.Println("  /plan     - (Stub) Coordinate planning agent")
			fmt.Println("  /code     - (Stub) Coordinate coding agent")
			fmt.Println("  /auto     - (Stub) Start Auto Drive orchestration")
			fmt.Println("  quit      - Quit the application")
		case "plan":
			fmt.Println("[Agent: Planner] Acknowledged. Ready to plan task.")
		case "code":
			fmt.Println("[Agent: Coder] Acknowledged. Ready to write code.")
		case "auto":
			fmt.Println("[Orchestrator: Auto Drive] Starting autonomous loop (Mock).")
		default:
			fmt.Printf("Unknown command: /%s\n", cmd)
		}
	} else {
		fmt.Printf("[Agent] Echoing input: %s\n", input)
	}
}
