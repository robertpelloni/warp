package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"github.com/robertpelloni/warp-rebuild/pkg/harness"
	"github.com/robertpelloni/warp-rebuild/pkg/agent"
)

type AgentService struct {
	server   *http.Server
	active   bool
	mu       sync.Mutex
	Provider harness.Provider
}

func NewAgentService(port int, provider harness.Provider) *AgentService {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "OK") })
	mux.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Prompt string `json:"prompt"` }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		a := agent.NewAgent("svc", "You are a helpful assistant.", &harness.ModelInfo{ID: "current"}, provider)
		a.History = append(a.History, harness.Message{Role: "user", Content: req.Prompt})
		if err := a.RunLoop(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(struct {
			History []harness.Message `json:"history"`
		}{History: a.History})
	})
	return &AgentService{
		server:   &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux},
		Provider: provider,
	}
}

func (s *AgentService) Start(ctx context.Context) error { s.mu.Lock(); s.active = true; s.mu.Unlock(); return s.server.ListenAndServe() }
func (s *AgentService) Stop(ctx context.Context) error { s.mu.Lock(); s.active = false; s.mu.Unlock(); return s.server.Shutdown(ctx) }
func (s *AgentService) IsActive() bool { s.mu.Lock(); defer s.mu.Unlock(); return s.active }
