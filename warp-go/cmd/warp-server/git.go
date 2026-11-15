package main

import (
	"encoding/json"
	"net/http"
	"log"
	"os/exec"
)

type GitRequest struct {
	Args []string `json:"args"`
}

type GitResponse struct {
	Output string `json:"output"`
}

func handleGit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[Git] Running command: git %v", req.Args)

	cmd := exec.Command("git", req.Args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[Git] Error: %v", err)
	}

	resp := GitResponse{
		Output: string(out),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
