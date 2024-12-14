package websocket

func main() {
	// Create a new websocket handler
	wh := NewWebsocketHandler()

	// Start the websocket server
	wh.Start()
}
