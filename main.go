package main

import (
	"log"
	"cco_api/database"
	"cco_api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database.InitDatabase()

	// Create Gin router
	r := gin.Default()

	// Register routes
	routes.RegisterRoutes(r)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
