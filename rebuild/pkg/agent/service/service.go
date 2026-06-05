package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/robertpelloni/warp-rebuild/pkg/agent"
	"github.com/robertpelloni/warp-rebuild/pkg/harness"
)

// AgentService manages AI agents and their sessions.
type AgentService struct {
	server *http.Server
	mu     sync.Mutex
	active bool
}

type RunRequest struct {
	Prompt string `json:"prompt"`
}

type RunResponse struct {
	History []harness.Message `json:"history"`
	Error   string            `json:"error,omitempty"`
}

func NewAgentService(port int) *AgentService {
	svc := &AgentService{}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	mux.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		var req RunRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Setup a minimal agent for the request
		model := &harness.ModelInfo{ID: "integrated-model"}
		mockProvider := &harness.MockProvider{
			Responses: []*harness.LLMResponse{
				{Content: fmt.Sprintf("Acknowledged: %s", req.Prompt)},
			},
		}
		a := agent.NewAgent("svc-agent", "You are an integrated Warp assistant.", model, mockProvider)
		a.History = append(a.History, harness.Message{Role: "user", Content: req.Prompt})

		err := a.RunLoop(context.Background())
		resp := RunResponse{History: a.History}
		if err != nil {
			resp.Error = err.Error()
		}

		json.NewEncoder(w).Encode(resp)
	})

	svc.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	return svc
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
