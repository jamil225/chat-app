package mywebsocket

import (
	"encoding/json"
	"errors"
	"log"
	"websocket/internal/httpclient"
	"websocket/internal/model"
)

//The WebSocketManager struct manages WebSocket connections,
//including registering and unregistering clients,
//and handling private messages between clients.

type WebSocketManager struct {
	Clients        map[int64]*Client   // Map of client ID to Client
	register       chan *Client        // Channel for registering new Clients
	unregister     chan *Client        // Channel for unregistering Clients
	privateMessage chan *model.Message // Channel for private messages
}

func NewHub() *WebSocketManager {
	return &WebSocketManager{
		Clients:        make(map[int64]*Client),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		privateMessage: make(chan *model.Message),
	}
}

// Run starts the main loop of the WebSocketManager, handling client registration,
// un-registration, and private messages
func (w *WebSocketManager) Run() {
	log.Println("WebSocketManager started!")
	for {
		select {
		case client := <-w.register:
			w.Clients[client.ID] = client
			callOnConnectApi(client.ID)
			log.Printf("Client connected: %s", client.ID)
		case client := <-w.unregister:
			if _, ok := w.Clients[client.ID]; ok {
				delete(w.Clients, client.ID)
				close(client.send)
				callOnDisConnectApi(client.ID)
				log.Printf("Client disconnected: %s", client.ID)
			}
		case message := <-w.privateMessage:
			if targetClient, ok := w.Clients[message.RecipientID]; ok {
				messageBytes, err := json.Marshal(message)
				if err != nil {
					log.Printf("Failed to marshal message: %v", err)
				}
				targetClient.send <- []byte(messageBytes)

				//TODO: send Message to reciver
				log.Printf("Message sent to client: %s", message.RecipientID)
			} else {
				//TODO: call DB to save message
				log.Printf("Recipient %d not connected", message.RecipientID)
			}
		}
	}
}

func (w *WebSocketManager) IsClientConnected(clientID int64) bool {
	_, exists := w.Clients[clientID]
	return exists
}

// GetClientByID retrieves a client by its ID from the WebSocketManager.
// Returns an error if the client is not found.
func (w *WebSocketManager) GetClientByID(clientID int64) (*Client, error) {
	log.Printf("Getting client by ID: %d", clientID)
	client, exists := w.Clients[clientID]
	if !exists {
		return nil, errors.New("client not found")
	}
	log.Printf("Client is conneted ID: %d", clientID)
	return client, nil
}

func callOnConnectApi(clientID int64) {
	log.Println("Calling chat-service onConnect")
	client := httpclient.NewHttpClient(0)    // 10 seconds timeout
	url := "http://localhost:9091/onConnect" // Include the protocol scheme
	postData := map[string]interface{}{
		"user_id": clientID,
		"topic":   "test-topic",
	}
	resp, err := client.Post(url, postData, map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		log.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()
}
func callOnDisConnectApi(clientID int64) {
	log.Println("Calling chat-service DisConnectApi")
	client := httpclient.NewHttpClient(0)       // 10 seconds timeout
	url := "http://localhost:9091/onDisConnect" // Include the protocol scheme
	postData := map[string]interface{}{
		"user_id": clientID,
		"topic":   "test-topic",
	}
	resp, err := client.Post(url, postData, map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		log.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()
}
