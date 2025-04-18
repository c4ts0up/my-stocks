package fetcher

import (
	"github.com/c4ts0up/my-stocks/backend/models"
	"strconv"
	"strings"
	"time"
)

// convertStockRatingsApiResponse converts StockRatingResponse to StockRating
func convertStockRatingsApiResponse(resp models.StockRatingRaw) (models.StockRating, error) {
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
		Action:     resp.Action,
		Brokerage:  resp.Brokerage,
		RatingFrom: resp.RatingFrom,
		RatingTo:   resp.RatingTo,
		Time:       parsedTime,
	}, nil
}

// parseDollarValue converts "$4.20" -> 4.20. It follows USA's money convention (, for 000's, . dor decimals)
func parseDollarValue(value string) (float64, error) {
	cleanValue := strings.ReplaceAll(value, "$", "")
	cleanValue = strings.ReplaceAll(cleanValue, ",", "")
	return strconv.ParseFloat(cleanValue, 64)
}
