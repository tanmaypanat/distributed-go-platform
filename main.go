package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default() // Creates a router with logging & recovery middleware

	// Define /orders endpoint
	router.GET("/orders", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Tanmay’s Go Distributed Platform!",
			"orders":  []string{"order1", "order2", "order3"},
		})
	})

	router.Run(":8080") // Start server on port 8080
}
