package client

import (
	"context"
	"time"

	pd "github.com/tanmaypanat/distributed-go-platform/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderClient struct {
	conn   *grpc.ClientConn
	client pd.OrderServiceClient
}

// NewOrderClient creates and returns a gRPC client wrapper
func NewOrderClient(address string) (*OrderClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := pd.NewOrderServiceClient(conn)
	return &OrderClient{
		conn:   conn,
		client: c,
	}, nil
}

// Close shuts down the gRPC connection
func (oc *OrderClient) Close() {
	if oc.conn != nil {
		oc.conn.Close()
	}
}

// GetOrder calls the gRPC server
func (oc *OrderClient) GetOrder(id string) (*pd.GetOrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return oc.client.GetOrder(ctx, &pd.GetOrderRequest{Id: id})
}
