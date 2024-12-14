package main

import (
	"github.com/IBM/sarama"
	"log"
)

func main() {
	producer, err := sarama.NewSyncProducer([]string{"192.168.1.5:9092"}, nil)
	if err != nil {
		log.Fatalf("Failed to start producer: %v", err)
	}
	defer func(producer sarama.SyncProducer) {
		err := producer.Close()
		if err != nil {
			log.Printf("Failed to close producer: %v", err)
		}
	}(producer)

	message := &sarama.ProducerMessage{Topic: "test-topic", Value: sarama.StringEncoder("Hello Kafka")}
	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
	log.Printf("Message sent to partition %d at offset %d", partition, offset)
}
