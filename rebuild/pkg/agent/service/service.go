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
	server *http.Server
	active bool
	mu sync.Mutex
}

func NewAgentService(port int) *AgentService {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "OK") })
	mux.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		var req struct { Prompt string `json:"prompt"` }
		json.NewDecoder(r.Body).Decode(&req)
		a := agent.NewAgent("svc", "prompt", &harness.ModelInfo{ID: "m"}, &harness.MockProvider{Responses: []*harness.LLMResponse{{Content: "resp"}}})
		a.RunLoop(context.Background())
		json.NewEncoder(w).Encode(struct{History []harness.Message `json:"history"`}{History: a.History})
	})
	return &AgentService{server: &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux}}
}

func (s *AgentService) Start(ctx context.Context) error { s.mu.Lock(); s.active = true; s.mu.Unlock(); return s.server.ListenAndServe() }
func (s *AgentService) Stop(ctx context.Context) error { s.mu.Lock(); s.active = false; s.mu.Unlock(); return s.server.Shutdown(ctx) }
func (s *AgentService) IsActive() bool { s.mu.Lock(); defer s.mu.Unlock(); return s.active }
