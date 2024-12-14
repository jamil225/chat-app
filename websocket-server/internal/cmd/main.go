package main

import (
	"log"
	"net/http"
	"websocket/internal/mysocket"
)

func main() {
	// Initialize the Hub
	hub := mysocket.NewHub()
	go hub.Run() // Start the Hub goroutine

	// Define the WebSocket route
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		mysocket.ServeWebSocket(hub, w, r)
	})

	// Start the HTTP server
	addr := ":8080"
	log.Printf("WebSocket server started on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
