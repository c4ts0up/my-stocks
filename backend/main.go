package main

import (
	"github.com/c4ts0up/my-stocks/backend/fetcher"
	"github.com/c4ts0up/my-stocks/backend/models"
	"github.com/c4ts0up/my-stocks/backend/presenter"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
	"time"
)

// Periodic fetcher (runs in the background)
func startPeriodicFetchAll(interval int, baseUrl string, api fetcher.IFetcherApi) {
	intervalDuration := time.Duration(interval) * time.Second
	ticker := time.NewTicker(intervalDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("🔄 Fetching data...")
			_ = api.FetchAll(baseUrl)
			log.Printf("✅ Data fetched")
		}
	}
}

func main() {
	log.Printf("--- MyStocks v0.0.2 ---")

	// Load environment variables (optional)
	dsn := os.Getenv("DATABASE_URL")
	fetchDelayStr := os.Getenv("FETCH_DELAY_S")
	apiUrl := os.Getenv("API_URL")

	fetchDelaySeconds, err := strconv.Atoi(fetchDelayStr)
	if err != nil {
		log.Fatalf("Failed to convert FETCH_DELAY to int: %v", err)
	}
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

	apiFetcher := fetcher.BasicStockRatingsFetcher{DB: models.DB, BearerToken: os.Getenv("API_BEARER_TOKEN")}
	go startPeriodicFetchAll(fetchDelaySeconds, apiUrl, &apiFetcher)

	// Set up the Gin router
	router := gin.Default()

	// Define routes
	router.GET("/stocks", presenter.GetStocks)
	router.GET("/stocks/:ticker", presenter.GetStockDetail)

	// Start the server
	log.Println("Server running at 0.0.0.0:8080")
	if err := router.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
