package mywebsocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"websocket/internal/httpclient"
	"websocket/internal/model"
	_ "websocket/internal/model"
)

// The Client struct represents a WebSocket client connected to the WebSocketManager.
type Client struct {
	webSocketManager *WebSocketManager
	conn             *websocket.Conn
	send             chan []byte
	ID               int64 // Unique client ID
}

// HandleReadMessages handles reading messages from the WebSocket connection
// and forwarding them to the WebSocketManager.
func (c *Client) HandleReadMessages() {
	log.Printf("Client connected: %s", c.ID)

	// Set a custom close callback handler
	c.conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("Client %d disconnected with code %d and reason: %s", c.ID, code, text)
		c.webSocketManager.unregister <- c
		return nil
	})
	defer func() {
		c.webSocketManager.unregister <- c // TODO: investigate this line
		err := c.conn.Close()
		if err != nil {
			return
		}
	}()
	log.Printf("Reading message from client %d", c.ID)
	for {
		_, rawMessage, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client %s: %v", c.ID, err)
			break
		}

		message := &model.Message{}
		if err := json.Unmarshal(rawMessage, message); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		//TODO: call chat api receive message
		callSendApi(c.ID, message)

		log.Printf("Received message from client %s: %s", c.ID, message.Content)
		//c.webSocketManager.privateMessage <- message // TODO: logic to call chat api receive message is it needed if I send alsway for reciver from my kafakf
	}
}

// TOdo: shit into httpclient
func callSendApi(id int64, message *model.Message) {
	log.Println("Calling chat-service sendMessage")
	client := httpclient.NewHttpClient(0)      // 10 seconds timeout
	url := "http://localhost:9091/sendMessage" // Include the protocol scheme
	postData := map[string]interface{}{
		"message":   message,
		"topic":     "test-topic",
		"partition": id,
	}
	resp, err := client.Post(url, postData, map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		log.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	var postResult map[string]interface{}
	err = httpclient.ParseResponseBody(resp, &postResult)
	if err != nil {
		log.Fatalf("Failed to parse response body: %v", err)
	}
	log.Printf("POST Response: %v", postResult)
}

// HandleReceivedMessage sends a received message to the client's send channel.
func (c *Client) HandleReceivedMessage(message model.Message) {
	log.Printf("Handling received message for client %d", c.ID)
	//TODO: investigate this which one to use
	//c.send <- message??
	//c.webSocketManager.privateMessage <- &message ??

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshalling message for client %d: %v", c.ID, err)
		return
	}

	err = c.conn.WriteMessage(websocket.TextMessage, messageBytes)
	if err != nil {
		log.Printf("Error sending message to client %d: %v", c.ID, err)
	}
}

// WriteMessagesToClient handles writing messages to the WebSocket connection
// from the send channel.
func (c *Client) WriteMessagesToClient() {
	defer func() {
		c.conn.Close()
	}()
	//TODO: call chat api send message
	log.Printf("Writing message to client %d", c.ID)
	for message := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Error sending message to client %s: %v", c.ID, err)
			break
		}
	}
}
