package agent

import (
	"fmt"
	"testing"
	"time"
)

func TestCircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker(2, 100*time.Millisecond)

	// First failure
	cb.Execute(func() error { return fmt.Errorf("fail") })
	if cb.GetState() != CircuitClosed {
		t.Error("Expected state Closed")
	}

	// Second failure - should trip
	cb.Execute(func() error { return fmt.Errorf("fail") })
	if cb.GetState() != CircuitOpen {
		t.Error("Expected state Open")
	}

	// Attempt while open
	err := cb.Execute(func() error { return nil })
	if err == nil || err.Error() != "circuit breaker is open" {
		t.Error("Expected error 'circuit breaker is open'")
	}

	// Wait for reset timeout
	time.Sleep(150 * time.Millisecond)
	cb.Execute(func() error { return nil })
	if cb.GetState() != CircuitClosed {
		t.Error("Expected state Closed after recovery")
	}
}
