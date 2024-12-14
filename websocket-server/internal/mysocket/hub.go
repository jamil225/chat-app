package mysocket

import "log"

type Hub struct {
	clients    map[int64]*Client // Map of client ID to Client
	register   chan *Client      // Channel for registering new clients
	unregister chan *Client      // Channel for unregistering clients
	private    chan *Message     // Channel for private messages
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int64]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		private:    make(chan *Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.ID] = client
			log.Printf("Client connected: %s", client.ID)
		case client := <-h.unregister:
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.send)
				log.Printf("Client disconnected: %s", client.ID)
			}
		case message := <-h.private:
			if targetClient, ok := h.clients[message.RecipientID]; ok {
				targetClient.send <- []byte(message.Content)
				log.Printf("Message sent to client: %s", message.RecipientID)
			} else {
				log.Printf("Recipient %d not connected", message.RecipientID)
			}
		}
	}
}

// Check if a client is connected by ID
func (h *Hub) IsClientConnected(clientID int64) bool {
	_, exists := h.clients[clientID]
	return exists
}
