package service

import (
	"context"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	s := NewAgentService(10005)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	go s.Start(ctx)
	time.Sleep(100 * time.Millisecond)
	if !s.IsActive() { t.Error("not active") }
	s.Stop(ctx)
}
