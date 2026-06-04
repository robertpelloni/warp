package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

// AgentService manages AI agents and their sessions.
type AgentService struct {
	server *http.Server
	mu     sync.Mutex
	active bool
}

func NewAgentService(port int) *AgentService {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	return &AgentService{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}
}

func (s *AgentService) Start(ctx context.Context) error {
	s.mu.Lock()
	s.active = true
	s.mu.Unlock()

	fmt.Printf("Agent service starting on %s\n", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *AgentService) Stop(ctx context.Context) error {
	s.mu.Lock()
	s.active = false
	s.mu.Unlock()
	return s.server.Shutdown(ctx)
}

func (s *AgentService) IsActive() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.active
}
