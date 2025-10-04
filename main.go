package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tanmaypanat/distributed-go-platform/client"
)

func main() {
	router := gin.Default() // Creates a router with logging & recovery middleware

	grpc_client, err := client.NewOrderClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer grpc_client.Close()
	// Define /orders endpoint
	router.GET("/orders/:id", func(c *gin.Context) {

		order_id := c.Param("id")

		if order_id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
			return
		}

		// Call gRPC server to get order details
		resp, err := grpc_client.GetOrder(order_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":          resp.GetId(),
			"description": resp.GetDescription(),
		})
	})

	log.Print("Starting HTTP server on :8080")
	router.Run(":8080") // Start server on port 8080
}
