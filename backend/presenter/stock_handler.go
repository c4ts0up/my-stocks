package presenter

import (
	"fmt"
	"net/http"

	"github.com/c4ts0up/my-stocks/backend/models"
	"github.com/c4ts0up/my-stocks/backend/presenter/models"
	"github.com/gin-gonic/gin"
)

// ConvertStockRating converts a DB model to a response model
func ConvertStockRating(dbModel models.StockRating) presenter.StockRatingResponse {
	return presenter.StockRatingResponse{
		Ticker:         dbModel.Ticker,
		TargetFrom:     formatFloat(dbModel.TargetFrom),
		TargetTo:       formatFloat(dbModel.TargetTo),
		Company:        dbModel.Company,
		Action:         dbModel.Action,
		Brokerage:      dbModel.Brokerage,
		RatingFrom:     dbModel.RatingFrom,
		RatingTo:       dbModel.RatingTo,
		Time:           dbModel.Time.Format("2006-01-02T15:04:05Z"),
		Recommendation: "Buy",
	}
}

// Helper to format floats cleanly
func formatFloat(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

// GetStockRatings handles the GET /stocks endpoint
func GetStockRatings(c *gin.Context) {
	var stockRatings []models.StockRating

	// Query database for all stock ratings
	if result := models.DB.Find(&stockRatings); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch stock ratings"})
		return
	}

	// Convert from DB models to response models
	var responseRatings []presenter.StockRatingResponse
	for _, stock := range stockRatings {
		responseRatings = append(responseRatings, ConvertStockRating(stock))
	}

	// Return the final response
	c.JSON(http.StatusOK, gin.H{"stock_ratings": responseRatings})
}
