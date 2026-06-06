package agent

import (
	"errors"
	"testing"
	"time"
)

func TestCircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker(2, 100*time.Millisecond)

	// Success
	err := cb.Execute(func() error { return nil })
	if err != nil { t.Errorf("expected nil, got %v", err) }
	if cb.GetState() != CircuitClosed { t.Errorf("expected closed, got %v", cb.GetState()) }

	// Failure 1
	cb.Execute(func() error { return errors.New("fail") })
	if cb.GetState() != CircuitClosed { t.Errorf("expected closed, got %v", cb.GetState()) }

	// Failure 2 (Threshold reached)
	cb.Execute(func() error { return errors.New("fail") })
	if cb.GetState() != CircuitOpen { t.Errorf("expected open, got %v", cb.GetState()) }

	// Immediate execution should fail
	err = cb.Execute(func() error { return nil })
	if err != ErrCircuitOpen { t.Errorf("expected ErrCircuitOpen, got %v", err) }

	// Wait for timeout
	time.Sleep(150 * time.Millisecond)
	if cb.GetState() != CircuitHalfOpen { t.Errorf("expected half-open, got %v", cb.GetState()) }

	// Success in half-open should close
	err = cb.Execute(func() error { return nil })
	if err != nil { t.Errorf("expected nil, got %v", err) }
	if cb.GetState() != CircuitClosed { t.Errorf("expected closed, got %v", cb.GetState()) }
}
