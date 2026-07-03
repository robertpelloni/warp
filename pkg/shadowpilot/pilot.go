package shadowpilot

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Pilot runs as a background process monitoring repository state.
type Pilot struct {
	mu          sync.Mutex
	lastDiff    string
	onAnomaly   func(message string)
	pollSeconds int
}

// New creates a new Shadow Pilot instance.
func New(onAnomaly func(string)) *Pilot {
	return &Pilot{
		onAnomaly:   onAnomaly,
		pollSeconds: 30,
	}
}

// Start begins the background monitoring loop.
func (p *Pilot) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(time.Duration(p.pollSeconds) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				p.checkState()
			}
		}
	}()
}

func (p *Pilot) checkState() {
	// Execute git diff to find uncommitted changes
	out, err := exec.Command("git", "diff").CombinedOutput()
	if err != nil {
		return // not a git repo or error running git
	}

	diffStr := string(out)
	if diffStr == "" {
		p.mu.Lock()
		p.lastDiff = ""
		p.mu.Unlock()
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// If the diff changed substantially, run a quick mock "anomaly" check
	if diffStr != p.lastDiff {
		p.lastDiff = diffStr

		// Mock logic: look for dangerous patterns in the diff
		if strings.Contains(diffStr, "TODO: remove") || strings.Contains(diffStr, "panic(") {
			if p.onAnomaly != nil {
				p.onAnomaly("[Shadow Pilot] Anomaly Detected: Potential debug code or unhandled panic spotted in uncommitted changes.")
			}
		}

		// Telemetry mock
		if strings.Contains(diffStr, "fmt.Println") {
			if p.onAnomaly != nil {
				p.onAnomaly(fmt.Sprintf("[Shadow Pilot] Info: Noticed %d line changes. Consider writing tests for new logic.", len(strings.Split(diffStr, "\n"))))
			}
		}
	}
}
