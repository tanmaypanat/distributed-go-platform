package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/tanmaypanat/distributed-go-platform/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedOrderServiceServer
}

func (s *server) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	orderID := req.GetId()

	var name string
	var price float64

	err := pgx_db.QueryRow(ctx, "SELECT name, price FROM orders WHERE id=$1", orderID).Scan(&name, &price)
	if err != nil {
		log.Printf("Error fetching order from DB: %v", err)
		return nil, err
	}

	description := fmt.Sprintf("Product name is %s -  It costs $%.2f $", name, price)
	// Create a response channel for this request
	respCh := make(chan *pb.GetOrderResponse)

	pendingRequests.Lock()
	pendingRequests.requests[orderID] = respCh
	pendingRequests.Unlock()

	// Produce to Kafka
	produceOrder(orderID, description)

	// Wait for response from Kafka consumer
	resp := <-respCh
	return resp, nil
}

func main() {
	initKafka()
	initPostgres()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, &server{})

	log.Println("Starting gRPC server on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
