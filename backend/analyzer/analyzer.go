package analyzer

import (
	"github.com/c4ts0up/my-stocks/backend/models"
	"time"
)

// IAnalyzerPipeline defines the pipeline interface
// It takes a Stock, a StockRating, and a list of AnalysisSteps
// and performs analysis.
type IAnalyzerPipeline interface {
	Analyze(stock *models.Stock, stockRating *models.StockRating, steps []IAnalysisStep)
}

// IAnalysisStep represents a single step in the analysis pipeline
type IAnalysisStep interface {
	Analyze(stock *models.Stock)
}

// BasicAnalyzerPipeline is a concrete implementation of IAnalyzerPipeline
type BasicAnalyzerPipeline struct {
	Steps []IAnalysisStep // not really a fan of this bit. If it's already an attribute, why am I passing it in Analyze?
}

func (b BasicAnalyzerPipeline) Analyze(stock *models.Stock, stockRating *models.StockRating, steps []IAnalysisStep) {
	for _, step := range steps {
		step.Analyze(stock)
	}
}

// DropStaleRecommendations drops stock recommendations if all the stock ratings happened over 3 months ago
type DropStaleRecommendations struct{}

func (m DropStaleRecommendations) Analyze(stock *models.Stock) {
	// Define the cutoff for stale ratings (3 months ago)
	staleCutoff := time.Now().AddDate(0, -3, 0)

	// Fetch all stock ratings for this stock
	var stockRatings []models.StockRating
	models.DB.Table("stock_ratings").Where("ticker = ?", stock.Ticker).Find(&stockRatings)

	// Track if any rating is recent
	hasRecentRating := false
	for _, rating := range stockRatings {
		if rating.Time.After(staleCutoff) {
			hasRecentRating = true
			break
		}
	}

	// If no recent rating is found, mark recommendation as N/A
	if !hasRecentRating {
		models.DB.Model(&models.Stock{}).
			Where("ticker = ?", stock.Ticker).
			Update("recommendation", "N/A")
	}
}
