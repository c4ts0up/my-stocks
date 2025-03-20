package fetcher

import (
	"github.com/c4ts0up/my-stocks/backend/models"
	"strconv"
	"strings"
	"time"
)

// convertApiResponse converts StockRatingResponse to StockRating
func convertApiResponse(resp models.StockRatingRaw) (models.StockRating, error) {
	// Strip dollar signs and convert to float
	targetFrom, err := parseDollarValue(resp.TargetFrom)
	if err != nil {
		return models.StockRating{}, err
	}

	targetTo, err := parseDollarValue(resp.TargetTo)
	if err != nil {
		return models.StockRating{}, err
	}

	// Parse time from string to time.Time
	parsedTime, err := time.Parse(time.RFC3339Nano, resp.Time)
	if err != nil {
		return models.StockRating{}, err
	}

	return models.StockRating{
		Ticker:     resp.Ticker,
		TargetFrom: targetFrom,
		TargetTo:   targetTo,
		Company:    resp.Company,
		Action:     resp.Action,
		Brokerage:  resp.Brokerage,
		RatingFrom: resp.RatingFrom,
		RatingTo:   resp.RatingTo,
		Time:       parsedTime,
	}, nil
}

// parseDollarValue converts "$4.20" -> 4.20
func parseDollarValue(value string) (float64, error) {
	cleanValue := strings.ReplaceAll(value, "$", "")
	return strconv.ParseFloat(cleanValue, 64)
}
