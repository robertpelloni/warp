package main

import (
	"encoding/json"
	"net/http"
	"log"
)

type AgentRequest struct {
	Message string `json:"message"`
}

type AgentResponse struct {
	Reply string `json:"reply"`
}

func handleAgent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AgentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[Agent] Received message: %s", req.Message)

	resp := AgentResponse{
		Reply: "I am Warp Go Agent. I received your message: " + req.Message,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
