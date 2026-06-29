package main

import (
	"encoding/json"
	"net/http"
)

type NewStatusResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Message string `json:"message"`
	Ready   bool   `json:"ready"`
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := NewStatusResponse{
		Status:  "online",
		Version: "1.0.0",
		Message: "Warp Agent HTTP/WebSocket Orchestrator Active",
		Ready:   true,
	}

	json.NewEncoder(w).Encode(resp)
}
