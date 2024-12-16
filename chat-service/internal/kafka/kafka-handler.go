package KafkaHandler

import (
	"chat-service.mod/internal/httpclient"
	"chat-service.mod/internal/model"
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
)

type KafkaHandler struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
}

func NewKafkaHandler(brokers []string) (*KafkaHandler, error) {
	log.Printf("----------Creating Kafka handler instance-------------")

	// Configure producer with manual partitioning
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.Partitioner = sarama.NewManualPartitioner
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Retry.Max = 5
	producerConfig.Producer.Return.Successes = true // Required for SyncProducer

	log.Printf("Creating Kafka producer!")
	producer, err := sarama.NewSyncProducer(brokers, producerConfig)
	if err != nil {
		log.Printf("Failed to start producer: %v", err)
		return nil, err
	}

	log.Printf("Creating Kafka consumer!")
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		err := producer.Close()
		if err != nil {
			log.Printf("Failed to close producer: %v", err)
			return nil, err
		}
		log.Printf("Failed to create new consumer: %v", err)
		return nil, err
	}

	log.Printf("Kafka handler created successfully")
	return &KafkaHandler{
		producer: producer,
		consumer: consumer,
	}, nil
}

func (k *KafkaHandler) StartProducer() error {
	return nil
}

func (k *KafkaHandler) StartConsumer(topic string, partition int32) error {
	log.Printf("Starting consumer for topic: %s, partition: %d", topic, partition)
	partitionConsumer, err := k.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Failed to start partition consumer: %v", err)
		return err
	}

	go func() {
		defer func() {
			log.Printf("Closing partition consumer")
			if err := partitionConsumer.Close(); err != nil {
				log.Printf("Failed to close partition consumer: %v", err)
			}
		}()
		for msg := range partitionConsumer.Messages() {
			k.onConsume(msg)
		}
	}()

	return nil
}

func (k *KafkaHandler) ProduceMessage(topic string, partition int32, message model.Message) error {
	// Create a producer message with the specified partition
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: partition, // Explicitly set partition
		Value:     sarama.StringEncoder(messageBytes),
	}

	partition, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Message sent to partition %d at offset %d", partition, offset)
	return nil
}

func (k *KafkaHandler) onConsume(msg *sarama.ConsumerMessage) {
	log.Println("Calling chat-service sendMessage")
	client := httpclient.NewHttpClient(0)    // 10 seconds timeout
	url := "http://localhost:8080/onReceive" // Include the protocol scheme

	message := model.Message{}
	if err := json.Unmarshal(msg.Value, &message); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	postData := map[string]interface{}{
		"message": message,
	}
	resp, err := client.Post(url, postData, map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		log.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	//var postResult map[string]interface{}
	//err = httpclient.ParseResponseBody(resp, &postResult)
	//if err != nil {
	//	log.Fatalf("Failed to parse response body: %v", err)
	//}
	//log.Printf("POST Response: %v", postResult)
	log.Printf("Consumed message offset %d, partition %d: %s", msg.Offset, msg.Partition, string(msg.Value))

}

func (k *KafkaHandler) Close() {
	if err := k.producer.Close(); err != nil {
		log.Printf("Failed to close producer: %v", err)
	}
	if err := k.consumer.Close(); err != nil {
		log.Printf("Failed to close consumer: %v", err)
	}
	log.Printf("Kafka handler closed successfully")
}

func (k *KafkaHandler) CloseConsumer(topic string, partition int32) error {
	partitionConsumer, err := k.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Failed to get partition consumer: %v", err)
		return err
	}

	if err := partitionConsumer.Close(); err != nil {
		log.Printf("Failed to close partition consumer: %v", err)
		return err
	}

	log.Printf("Partition consumer for topic %s and partition %d closed successfully", topic, partition)
	return nil
}
