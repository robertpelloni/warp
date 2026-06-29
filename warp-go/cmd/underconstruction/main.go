package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type AgentState struct {
	LastActive string `json:"last_active"`
	RunCount   int    `json:"run_count"`
	Status     string `json:"status"`
}

func main() {
	cacheDir := ".pi"
	stateFile := filepath.Join(cacheDir, "agent_state.json")

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		log.Fatalf("Failed to create cache dir: %v", err)
	}

	state := AgentState{
		LastActive: time.Now().Format(time.RFC3339),
		RunCount:   1,
		Status:     "running",
	}

	if data, err := ioutil.ReadFile(stateFile); err == nil {
		var oldState AgentState
		if err := json.Unmarshal(data, &oldState); err == nil {
			state.RunCount = oldState.RunCount + 1
		}
	}

	fmt.Printf("Warp Core Agent Loop Initializing... (Run #%d)\n", state.RunCount)
	fmt.Println("Establishing persistent connection to main engine...")

	// Simulate work loop
	for i := 0; i < 3; i++ {
		fmt.Printf("Processing tick %d...\n", i)
		time.Sleep(500 * time.Millisecond)
	}

	state.Status = "completed"
	state.LastActive = time.Now().Format(time.RFC3339)

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		log.Fatalf("Failed to serialize state: %v", err)
	}

	if err := ioutil.WriteFile(stateFile, data, 0644); err != nil {
		log.Fatalf("Failed to write state: %v", err)
	}

	fmt.Println("Agent loop finished successfully. State persisted.")
}
