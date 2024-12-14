package main

import (
	"github.com/IBM/sarama"
	"log"
)

func main() {
	consumer, err := sarama.NewConsumer([]string{"192.168.1.5:9092"}, nil)
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("test-topic", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to start partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
		log.Printf("Consumed message offset %d: %s", msg.Offset, string(msg.Value))
		//call send message onConsume()
	}
}
