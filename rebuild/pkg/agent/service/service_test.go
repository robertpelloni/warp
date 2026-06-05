package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestAgentService(t *testing.T) {
	s := NewAgentService(10002)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		_ = s.Start(ctx)
	}()

	time.Sleep(200 * time.Millisecond)
	if !s.IsActive() {
		t.Error("Expected service to be active")
	}

	// Test /run endpoint
	reqBody, _ := json.Marshal(RunRequest{Prompt: "Hello Warp"})
	resp, err := http.Post("http://localhost:10002/run", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to call /run: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var runResp RunResponse
	if err := json.NewDecoder(resp.Body).Decode(&runResp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(runResp.History) < 3 { // system + user + assistant
		t.Errorf("Expected at least 3 messages in history, got %d", len(runResp.History))
	}

	_ = s.Stop(ctx)
	if s.IsActive() {
		t.Error("Expected service to be inactive")
	}
}
