package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	pb "github.com/tanmaypanat/distributed-go-platform/proto"
)

var kafkaWriter *kafka.Writer
var kafkaReader *kafka.Reader

// Store pending requests
var pendingRequests = struct {
	sync.Mutex
	requests map[string]chan *pb.GetOrderResponse
}{requests: make(map[string]chan *pb.GetOrderResponse)}

func initKafka() {
	kafkaWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{"localhost:9092"},
		Topic:        "orders",
		Balancer:     &kafka.LeastBytes{},
		Async:        false,
		BatchTimeout: 10 * time.Millisecond, // lower batching delay
		RequiredAcks: 1,                     // respond quicker
	})

	kafkaReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"},
		GroupID:        "order-service",
		Topic:          "orders",
		MinBytes:       1,
		MaxBytes:       10e6,
		MaxWait:        10 * time.Millisecond,
		SessionTimeout: 30 * time.Second,
	})

	go startKafkaConsumer()
}

func produceOrder(orderID, description string) {
	msg := kafka.Message{
		Key:   []byte(orderID),
		Value: []byte(description),
	}
	err := kafkaWriter.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Println("Error producing order:", err)
	} else {
		log.Println("Produced order:", orderID, description)
	}
}

func startKafkaConsumer() {
	for {
		msg, err := kafkaReader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading Kafka message:", err)
			continue
		}

		orderID := string(msg.Key)
		description := string(msg.Value)

		log.Printf("Consumed order: id=%s description=%s\n", orderID, description)

		// If there is a pending request for this order ID, send response
		pendingRequests.Lock()
		if ch, ok := pendingRequests.requests[orderID]; ok {
			ch <- &pb.GetOrderResponse{Id: orderID, Description: description}
			delete(pendingRequests.requests, orderID)
		}
		pendingRequests.Unlock()
	}
}
