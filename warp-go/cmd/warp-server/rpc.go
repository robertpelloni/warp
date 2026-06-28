package main

import (
	"encoding/json"
	"net/http"
	"log"
)

type RpcRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  json.RawMessage `json:"params"`
	Id      interface{} `json:"id"`
}

type RpcResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Id      interface{} `json:"id"`
}

func handleRPC(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RpcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[RPC] Received method: %s", req.Method)

	resp := RpcResponse{
		Jsonrpc: "2.0",
		Id:      req.Id,
		Result:  "Method " + req.Method + " executed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
