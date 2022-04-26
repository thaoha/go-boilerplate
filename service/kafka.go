package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type KafkaBroker struct {
	brokers  string
	producer *kafka.Producer
}

func (k *KafkaBroker) NewConsumer(topics []string, group string, offset string) *kafka.Consumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": k.brokers,
		"group.id": group,
		"auto.offset.reset": offset,
	})

	if err != nil {
		panic(err)
	}
	consumer.SubscribeTopics(topics, nil)
	return consumer
}

func (k *KafkaBroker) NewProducer() *kafka.Producer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": k.brokers,
	})
	if err != nil {
		panic(err)
	}
	log.Printf("Created producer %v\n", producer)
	return producer
}

func (k *KafkaBroker) Produce(topic string, data interface{}) {
	if k.producer == nil {
		k.producer = k.NewProducer()
	}
	log.Printf("Sending message to %v\n", topic)
	deliveryChan := make(chan kafka.Event)
	message, _ := json.Marshal(data)

	k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, deliveryChan)

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		log.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)
}

func (k *KafkaBroker) CreateTopics(topics []string) {
	a, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": k.brokers})
	if err != nil {
		log.Printf("Failed to create Admin client: %s\n", err)
		os.Exit(1)
	}

	// Contexts are used to abort or limit the amount of time
	// the Admin call blocks waiting for a result.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create topics on cluster.
	// Set Admin options to wait for the operation to finish (or at most 60s)
	maxDur, err := time.ParseDuration("60s")
	if err != nil {
		panic("ParseDuration(60s)")
	}
	var topicSpecifications []kafka.TopicSpecification
	for _, topic := range topics {
		topicSpecifications = append(topicSpecifications, kafka.TopicSpecification{
			Topic:             topic,
			NumPartitions:     5,
			ReplicationFactor: 2,
		})
	}
	results, err := a.CreateTopics(
		ctx,
		topicSpecifications,
		kafka.SetAdminOperationTimeout(maxDur))
	if err != nil {
		log.Printf("Failed to create topic: %v\n", err)
		os.Exit(1)
	}

	// Print results
	for _, result := range results {
		log.Printf("%s\n", result)
	}

	a.Close()
}