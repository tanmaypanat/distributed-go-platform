package main

import (
	"context"
	"log"
	"time"

	pd "github.com/tanmaypanat/distributed-go-platform/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pd.NewOrderServiceClient(conn)

	req := &pd.GetOrderRequest{Id: "145"}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetOrder(ctx, req)
	if err != nil {
		log.Fatalf("Error while calling GetOrder: %v", err)
	}

	log.Printf("Response from server. Id = %s , description = %s", res.GetId(), res.GetDescription())
}
