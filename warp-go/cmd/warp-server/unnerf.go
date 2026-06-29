package main

import (
	"encoding/json"
	"net/http"
	"log"

	"github.com/robertpelloni/warp/warp-go/cmd/warp-server/nerf"
)

type UnnerfRequest struct {
	Scene string `json:"scene"`
}

type UnnerfResponse struct {
	Result string `json:"result"`
}

func handleUnnerf(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UnnerfRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[Unnerf] Rendering scene: %s", req.Scene)

	// Invoke the neural rendering pipeline
	nerf.RenderScene(64, 64) // low res for speed demo

	resp := UnnerfResponse{
		Result: "Unnerf Neural Rendering pipeline complete for scene: " + req.Scene,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
