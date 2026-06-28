package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"log"
)

type FsRequest struct {
	Path string `json:"path"`
}

type FsResponse struct {
	Content string `json:"content"`
}

func handleFs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req FsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[FS] Reading file: %s", req.Path)

	content, err := ioutil.ReadFile(req.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := FsResponse{
		Content: string(content),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
