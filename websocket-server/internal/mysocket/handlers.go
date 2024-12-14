package mysocket

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
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
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
		ID:   clientID,
	}
	hub.register <- client

	go client.ReadPump()
	go client.WritePump()
}
