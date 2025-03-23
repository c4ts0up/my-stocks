package main

import (
	"github.com/c4ts0up/my-stocks/backend/analyzer"
	"github.com/c4ts0up/my-stocks/backend/fetcher"
	"github.com/c4ts0up/my-stocks/backend/models"
	"github.com/c4ts0up/my-stocks/backend/presenter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
	"time"
)

// Periodic fetcher (runs in the background)
func startPeriodicFetchAll(interval int, ratingsUrl string, infoUrl string, api *fetcher.StockFetcher) {
	intervalDuration := time.Duration(interval) * time.Second
	ticker := time.NewTicker(intervalDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("🔄 Fetching data...")
			_ = api.FetchAll(ratingsUrl, infoUrl)
			log.Printf("✅ Data fetched")
		}
	}
}

func startPeriodicAnalysis(interval int, analyzerPipeline analyzer.IAnalyzerPipeline) {
	intervalDuration := time.Duration(interval) * time.Second
	ticker := time.NewTicker(intervalDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("🔄 Getting all stocks ...")
			stocks, err := models.GetAllStocks()
			if err != nil {
				log.Printf("Couldn't fetch all stocks, %e\n", err)
			}

			for _, stock := range stocks {
				analyzerPipeline.Analyze(&stock)
			}
			log.Printf("✅ Data fetched")
		}
	}
}

func main() {
	log.Printf("--- MyStocks v0.0.2 ---")

	// Load environment variables (optional)
	dsn := os.Getenv("DATABASE_URL")
	fetchDelayStr := os.Getenv("FETCH_DELAY_S")
	analysisDelayStr := os.Getenv("ANALYSIS_DELAY_S")

	ratingsApiUrl := os.Getenv("RATINGS_API_URL")
	ratingsApiToken := os.Getenv("RATINGS_API_TOKEN")
	infoApiUrl := os.Getenv("INFO_API_URL")
	infoApiToken := os.Getenv("INFO_API_TOKEN")

	fetchDelaySeconds, err := strconv.Atoi(fetchDelayStr)
	if err != nil {
		log.Fatalf("Failed to convert FETCH_DELAY to int: %v", err)
	}
	analysisDelaySeconds, err := strconv.Atoi(analysisDelayStr)
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

	apiFetcher := fetcher.StockFetcher{
		RatingsFetcher: &fetcher.BasicStockRatingsFetcher{DB: models.DB, BearerToken: ratingsApiToken},
		InfoFetcher:    &fetcher.BasicStockInfoFetcher{DB: models.DB, BearerToken: infoApiToken},
	}

	analyzerPipeline := analyzer.BasicAnalyzerPipeline{}

	go startPeriodicFetchAll(fetchDelaySeconds, ratingsApiUrl, infoApiUrl, &apiFetcher)
	go startPeriodicAnalysis(analysisDelaySeconds, &analyzerPipeline)

	// Set up the Gin router
	router := gin.Default()

	// Enable CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("FRONTEND_URL")},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Content-Type,Authorization"},
		AllowCredentials: true,
	}))

	// Define routes
	router.GET("/stocks", presenter.GetStocks)
	router.GET("/stocks/:ticker", presenter.GetStockDetail)

	// Start the server
	log.Println("Server running at 0.0.0.0:8080")
	if err := router.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
