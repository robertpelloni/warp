package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[WebSocket] Upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("[WebSocket] Client connected")

	// Send an initial welcome message
	err = conn.WriteMessage(websocket.TextMessage, []byte(`{"type": "system", "message": "Connected to Warp Go Server"}`))
	if err != nil {
		log.Println("[WebSocket] Write error:", err)
		return
	}

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("[WebSocket] Read error:", err)
			break
		}

		log.Printf("[WebSocket] Received: %s\n", p)

		// Echo message back to client for now
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			log.Println("[WebSocket] Echo error:", err)
			break
		}
	}
	log.Println("[WebSocket] Client disconnected")
}
