package main

import (
	"encoding/json"
	"net/http"
	"log"
)

type CdpRequest struct {
	Url string `json:"url"`
}

type CdpResponse struct {
	Result string `json:"result"`
}

func handleCdp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CdpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[CDP] Inspecting network activity for: %s", req.Url)

	resp := CdpResponse{
		Result: "Network activity for " + req.Url + " logged successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
