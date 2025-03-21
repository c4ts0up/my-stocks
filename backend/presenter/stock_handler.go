package presenter

import (
	"fmt"
	"github.com/c4ts0up/my-stocks/backend/models"
	presenter "github.com/c4ts0up/my-stocks/backend/presenter/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetStocks handles GET /stocks
func GetStocks(c *gin.Context) {
	var stocks []models.Stock

	if result := models.DB.Find(&stocks); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch stocks"})
		return
	}

	stockBases := make([]presenter.StockBase, len(stocks))
	for i, s := range stocks {
		stockBases[i] = presenter.StockBase{
			Ticker:         s.Ticker,
			CompanyName:    s.Company,
			CurrentPrice:   "N/A", // Placeholder for now
			Recommendation: s.Recommendation,
		}
	}

	c.JSON(http.StatusOK, stockBases)
}

// GetStockDetail handles GET /stocks/:ticker
func GetStockDetail(c *gin.Context) {
	ticker := c.Param("ticker")

	var stock models.Stock
	if result := models.DB.Where("ticker = ?", ticker).First(&stock); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "stock not found"})
		return
	}

	var stockRatings []models.StockRating
	models.DB.Where("ticker = ?", ticker).Find(&stockRatings)

	ratings := make([]presenter.StockRating, len(stockRatings))
	for i, r := range stockRatings {
		ratings[i] = presenter.StockRating{
			TargetFrom: formatFloat(r.TargetFrom),
			TargetTo:   formatFloat(r.TargetTo),
			Action:     r.Action,
			Brokerage:  r.Brokerage,
			RatingFrom: r.RatingFrom,
			RatingTo:   r.RatingTo,
			Time:       r.Time.Format("2006-01-02T15:04:05Z"),
		}
	}

	response := presenter.StockDetail{
		StockBase: presenter.StockBase{
			Ticker:         stock.Ticker,
			CompanyName:    stock.Company,
			CurrentPrice:   "N/A", // Placeholder
			Recommendation: stock.Recommendation,
		},
		StockRatings: ratings,
	}

	c.JSON(http.StatusOK, response)
}

func formatFloat(value float64) string {
	return fmt.Sprintf("%.2f", value)
}
