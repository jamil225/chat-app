package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"websocket/internal/model"
	"websocket/internal/mywebsocket"
)

type RequestHandler struct {
	webSocketManager *mywebsocket.WebSocketManager
}

func NewRequestHandler(webSocketManager *mywebsocket.WebSocketManager) *RequestHandler {
	return &RequestHandler{webSocketManager: webSocketManager}
}

func (h *RequestHandler) OnReceiveMessage(c *gin.Context) {
	log.Default().Println("Receiving message for client")
	var request struct {
		Message model.Message `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := h.webSocketManager.GetClientByID(request.Message.RecipientID) // TODO : change if client vs partition id mapping changes
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		//  TODO: handle by saving it on DB
		return
	}
	log.Default().Println("Receiving message for client %d", client.ID)

	client.HandleReceivedMessage(request.Message)

	c.JSON(http.StatusOK, gin.H{"status": "message sent"})
}
