package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Warp Go Server Backend is operational.\n\nUse /status for health checks.")
	})

	mux.HandleFunc("/agent", handleAgent)
	mux.HandleFunc("/fs", handleFs)
	mux.HandleFunc("/git", handleGit)
	mux.HandleFunc("/ast", handleAst)
	mux.HandleFunc("/cdp", handleCdp)
	mux.HandleFunc("/nerf", handleNerf)
	mux.HandleFunc("/unnerf", handleUnnerf)
	mux.HandleFunc("/rpc", handleRPC)
	mux.HandleFunc("/ws", handleWebSocket)
	mux.HandleFunc("/status", handleStatus)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMiddleware(mux),
	}

	go func() {
		log.Printf("[Init] Warp Server starting on HTTP :%s...", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[Error] Could not listen on %s: %v\n", port, err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Println("\n[Shutdown] Shutting down server gracefully...")

	// Implementation for graceful shutdown with timeout context can be added here
	time.Sleep(1 * time.Second)
	log.Println("[Shutdown] Server stopped.")
}
