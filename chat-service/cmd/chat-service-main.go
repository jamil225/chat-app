package main

import (
	myRequestHandler "chat-service.mod/internal/api"
	myKafkaHandler "chat-service.mod/internal/kafka"
	"chat-service.mod/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
)

func main() {
	log.Println("Starting chat service!")
	kafkaHandler := setupAndStartKafkaHandler()
	defer kafkaHandler.Close()
	startHTTPServer(kafkaHandler)
}

func startHTTPServer(kafkaHandler *myKafkaHandler.KafkaHandler) {
	requestHandler := myRequestHandler.NewRequestHandler(kafkaHandler)
	router := gin.Default()
	router.POST("/sendMessage", requestHandler.SendMessage)
	router.POST("/onConnect", requestHandler.OnConnect)
	router.POST("/onDisConnect", requestHandler.OnDisConnect)

	if err := router.Run(":9091"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func setupAndStartKafkaHandler() *myKafkaHandler.KafkaHandler {
	log.Println("Initializing Kafka handler!")

	brokers := getConfigProperty("kafka.brokers")
	kafkaHandler, err := myKafkaHandler.NewKafkaHandler(brokers)
	if err != nil {
		log.Fatalf("Failed to create Kafka handler: %v", err)
	}

	if err := kafkaHandler.StartProducer(); err != nil {
		log.Fatalf("Failed to start producer: %v", err)
	}

	if err := kafkaHandler.StartConsumer("test-topic", 0); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}
	sampleMessage := model.Message{
		Content:     "Hello, on start ping from 1->2",
		SenderID:    1,
		RecipientID: 2,
	}
	if err := kafkaHandler.ProduceMessage("test-topic", 0, sampleMessage); err != nil {
		log.Printf("Failed to produce message: %v", err)
	}
	log.Println("Kafka handler initialized successfully!")

	return kafkaHandler
}

func getConfigProperty(key string) []string {
	log.Println("Reading config properties!")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs") // Relative path from the root folder `chat-service`
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	value := viper.GetStringSlice(key)
	return value
}
