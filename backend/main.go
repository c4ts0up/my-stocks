package main

import (
	"github.com/c4ts0up/my-stocks/backend/models"
	"github.com/c4ts0up/my-stocks/backend/presenter"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	log.Printf("--- MyStocks v0.0.2 ---")

	// Load environment variables (optional)
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgresql://root@localhost:26257/stocks_db?sslmode=disable"
	}

	// Connect to the database
	if err := models.ConnectDB(dsn); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer func() {
		if err := models.CloseDB(); err != nil {
			log.Printf("Error closing DB: %v", err)
		}
	}()

	// Set up the Gin router
	router := gin.Default()

	// Define routes
	router.GET("/stocks", presenter.GetStockRatings)

	// Start the server
	log.Println("Server running at 0.0.0.0:8080")
	if err := router.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
