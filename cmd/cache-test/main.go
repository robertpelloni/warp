package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/robertpelloni/warp/pkg/agent"
)

func main() {
	fmt.Println("Starting Cache Test Harness...")

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting CWD: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Workspace Root: %s\n", cwd)

	orch, err := agent.NewOrchestrator(cwd)
	if err != nil {
		fmt.Printf("Failed to initialize orchestrator: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Initial Status: %s\n", orch.Status())

	fmt.Println("Simulating Remote Connection...")
	if err := orch.ConnectRemote(); err != nil {
		fmt.Printf("Failed to connect remote: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Status After Connect: %s\n", orch.Status())

	// Verify the file was created
	piDir := filepath.Join(cwd, ".pi")
	stateFile := filepath.Join(piDir, "state.json")

	if _, err := os.Stat(stateFile); err == nil {
		fmt.Println("SUCCESS: .pi/state.json was successfully created and persisted.")
	} else {
		fmt.Printf("ERROR: Expected %s to exist, but it doesn't.\n", stateFile)
		os.Exit(1)
	}
}
