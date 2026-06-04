package service

import (
	"context"
	"testing"
	"time"
)

func TestAgentService(t *testing.T) {
	s := NewAgentService(9999)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go func() {
		_ = s.Start(ctx)
	}()

	time.Sleep(100 * time.Millisecond)
	if !s.IsActive() {
		t.Error("Expected service to be active")
	}

	_ = s.Stop(ctx)
	if s.IsActive() {
		t.Error("Expected service to be inactive")
	}
}
