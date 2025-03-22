package analyzer

import (
	"github.com/c4ts0up/my-stocks/backend/models"
	"slices"
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

// Analyze runs the hardcoded analysis steps on the given stock
func (b *BasicAnalyzerPipeline) Analyze(stock *models.Stock) {
	// Hardcoded steps in the desired order
	steps := []IAnalysisStep{
		PriceChangePonderedRecommendation{},
		DropStaleRecommendations{},
	}

	// Execute each step
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

// PriceChangePonderedRecommendation changes the recommendation if the price has
type PriceChangePonderedRecommendation struct{}

func (m PriceChangePonderedRecommendation) Analyze(stock *models.Stock) {
	// Creates the map with target keyword frequency
	targetFrequency := map[string]int{
		"Buy":  0,
		"Hold": 0,
		"Sell": 0,
	}

	positiveTargets := []string{"Buy", "Overweight", "Outperform", "Market Outperform"}
	neutralTargets := []string{"Neutral", "Equal Weight", "Perform", "Market Perform"}
	negativeTargets := []string{"Sell", "Underweight", "Underperform", "Market Underperform"}

	// Fetch all stock ratings for this stock
	var stockRatings []models.StockRating
	models.DB.Table("stock_ratings").Where("ticker = ?", stock.Ticker).Find(&stockRatings)

	// Maps positive, negative and neutral ratings to three categories
	for _, rating := range stockRatings {
		if slices.Contains(positiveTargets, rating.RatingTo) {
			targetFrequency["Buy"]++
		} else if slices.Contains(neutralTargets, rating.RatingTo) {
			targetFrequency["Hold"]++
		} else if slices.Contains(negativeTargets, rating.RatingTo) {
			targetFrequency["Sell"]++
		}
	}

	// Gets the best frequency
	// Solves the doubt "what if they're tied?" Hardcoded recommendation order. Could become a feature, though
	// like "investor personality"
	recommendationOrder := []string{"Hold", "Sell", "Buy"}

	maxFrequency := 0
	maxRecommendation := "N/A"

	for _, r := range recommendationOrder {
		if targetFrequency[r] > maxFrequency {
			maxFrequency = targetFrequency[r]
			maxRecommendation = r
		}
	}

	// Update the stock's recommendation
	stock.Recommendation = maxRecommendation

	// Save the updated stock to the database
	models.DB.Save(&stock)
}
