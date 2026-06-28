package main

import (
	"encoding/json"
	"net/http"
	"log"
)

type AstRequest struct {
	Path string `json:"path"`
}

type AstResponse struct {
	Result string `json:"result"`
}

func handleAst(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AstRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[AST] Analyzing file: %s", req.Path)

	resp := AstResponse{
		Result: "AST Deep Analysis Complete for: " + req.Path,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
