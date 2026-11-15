package main

import (
	"encoding/json"
	"net/http"
	"log"
)

type NerfRequest struct {
	Scene string `json:"scene"`
}

type NerfResponse struct {
	Result string `json:"result"`
}

func handleNerf(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req NerfRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[NeRF] Rendering scene: %s", req.Scene)

	resp := NerfResponse{
		Result: "Neural Radiance Field rendering complete for scene: " + req.Scene,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
