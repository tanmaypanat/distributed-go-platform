package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pd "github.com/tanmaypanat/distributed-go-platform/proto"
	"google.golang.org/grpc"
)

type server struct {
	pd.UnimplementedOrderServiceServer
}

func (s *server) GetOrder(ctx context.Context, req *pd.GetOrderRequest) (*pd.GetOrderResponse, error) {
	log.Printf("Recieved order request with id: %s", req.GetId())

	return &pd.GetOrderResponse{
		Id:          req.GetId(),
		Description: fmt.Sprintf("Order with id %s is a test order", req.GetId()),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pd.RegisterOrderServiceServer(grpcServer, &server{})

	log.Println("Starting gRPC server on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
