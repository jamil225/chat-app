package requestHandler

import (
	myKafkaHandler "chat-service.mod/internal/kafka"
	"chat-service.mod/internal/model"
	_ "chat-service.mod/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RequestHandler struct {
	kafkaHandler *myKafkaHandler.KafkaHandler
}

func NewRequestHandler(kafkaHandler *myKafkaHandler.KafkaHandler) *RequestHandler {
	return &RequestHandler{kafkaHandler: kafkaHandler}
}

func (h *RequestHandler) SendMessage(c *gin.Context) {
	var request struct {
		Message   model.Message `json:"message" binding:"required"`
		Topic     string        `json:"topic" binding:"required"`
		Partition int32         `json:"partition" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.kafkaHandler.ProduceMessage(request.Topic, request.Partition, request.Message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "message sent"})
}

func (h *RequestHandler) OnConnect(c *gin.Context) {
	var request struct {
		UserID int32  `json:"user_id" binding:"required"`
		Topic  string `json:"topic" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.kafkaHandler.StartConsumer(request.Topic, request.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "consumer already running!!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "consumer started"})
}

func (h *RequestHandler) OnDisConnect(c *gin.Context) {
	//TODO: closing all how to fix
	var request struct {
		UserID int32  `json:"user_id" binding:"required"`
		Topic  string `json:"topic" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.kafkaHandler.CloseConsumer(request.Topic, request.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "consumer already running!!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "consumer started"})
}
