package mywebsocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

/*ServeWebSocket handles upgrading an HTTP connection to a WebSocket connection,
*creates a new Client, registers it with the WebSocketManager, and starts
* goroutines to handle reading and writing messages for the WebSocket connection.
 */
func ServeWebSocket(webSocketManager *WebSocketManager, w http.ResponseWriter, r *http.Request) {
	clientIDStr := r.URL.Query().Get("id")
	if clientIDStr == "" {
		http.Error(w, "Missing client ID", http.StatusBadRequest)
		return
	}

	clientID, err := strconv.ParseInt(clientIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}

	client := &Client{
		webSocketManager: webSocketManager,
		conn:             conn,
		send:             make(chan []byte, 256),
		ID:               clientID,
	}
	// Register the client with the WebSocketManager
	webSocketManager.register <- client
	//callOnConnectApi(clientID) why??
	go client.HandleReadMessages()
	go client.WriteMessagesToClient()
}
