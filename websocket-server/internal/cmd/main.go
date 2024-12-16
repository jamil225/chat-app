package main

import (
	"github.com/gin-gonic/gin"
	"log"
	MyRequestHandler "websocket/internal/api"
	"websocket/internal/mywebsocket"
)

func main() {
	// Initialize the WebSocketManager
	hub := mywebsocket.NewHub()
	go hub.Run() // Start the WebSocketManager goroutine

	// Create a new Gin router
	router := gin.Default()

	// Define the WebSocket route
	router.GET("/ws", func(c *gin.Context) {
		mywebsocket.ServeWebSocket(hub, c.Writer, c.Request)
	})

	requestHandler := MyRequestHandler.NewRequestHandler(hub)
	router.POST("/onReceive", requestHandler.OnReceiveMessage)

	// Start the HTTP server on port 8080
	addr := ":8080"
	log.Printf("Server started on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
