package analyzer

import (
	"github.com/c4ts0up/my-stocks/backend/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDropStaleRecommendations_AllStale_DBUpdate(t *testing.T) {
	// Setup stock and stale ratings
	stock := models.Stock{Ticker: "AAPL", Recommendation: "BUY"}
	staleTime := time.Now().AddDate(0, -4, 0)
	stockRatings := []models.StockRating{
		{Ticker: "AAPL", Brokerage: "A", Time: staleTime},
		{Ticker: "AAPL", Brokerage: "B", Time: staleTime},
	}

	// Mock database call
	models.DB = models.NewTestDB(stockRatings)
	models.DB.Create(&stock) // Save stock to DB first

	// Run analysis
	analyzer := DropStaleRecommendations{}
	analyzer.Analyze(&stock)

	// Reload stock from DB to verify persistence
	var updatedStock models.Stock
	models.DB.First(&updatedStock, "ticker = ?", "AAPL")

	// Assert recommendation is dropped in DB
	assert.Equal(t, "N/A", updatedStock.Recommendation)
}

func TestDropStaleRecommendations_HasRecent(t *testing.T) {
	// Setup stock and a mix of stale and recent ratings
	stock := &models.Stock{Ticker: "GOOGL", Recommendation: "BUY"}
	staleTime := time.Now().AddDate(0, -4, 0)
	recentTime := time.Now().AddDate(0, -1, 0)
	stockRatings := []models.StockRating{
		{Ticker: "GOOGL", Brokerage: "A", Time: staleTime},
		{Ticker: "GOOGL", Brokerage: "B", Time: recentTime},
	}

	// Mock database call
	models.DB = models.NewTestDB(stockRatings)
	models.DB.Create(&stock)

	// Run analysis
	analyzer := DropStaleRecommendations{}
	analyzer.Analyze(stock)

	// Reload stock from DB and assert recommendation remains unchanged
	var updatedStock models.Stock
	models.DB.First(&updatedStock, "ticker = ?", stock.Ticker)
	assert.Equal(t, "BUY", updatedStock.Recommendation)
}

func TestDropStaleRecommendations_NoRatings(t *testing.T) {
	// Setup stock with no ratings
	stock := &models.Stock{Ticker: "MSFT", Recommendation: "HOLD"}
	stockRatings := []models.StockRating{}

	// Mock database call
	models.DB = models.NewTestDB(stockRatings)
	models.DB.Create(stock)

	// Run analysis
	analyzer := DropStaleRecommendations{}
	analyzer.Analyze(stock)

	// Reload stock from DB and assert recommendation is dropped
	var updatedStock models.Stock
	models.DB.First(&updatedStock, "ticker = ?", stock.Ticker)
	assert.Equal(t, "N/A", updatedStock.Recommendation)
}
