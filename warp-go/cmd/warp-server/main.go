package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type StatusResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Warp Go Server Backend is operational.\n\nUse /status for health checks.")
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		resp := StatusResponse{
			Status:  "online",
			Version: "0.1.1",
			Message: "Warp Agent HTTP/WebSocket Orchestrator Active",
		}

		json.NewEncoder(w).Encode(resp)
	})

	port := ":8080"
	fmt.Printf("Warp Server Initialized - Listening on HTTP %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
