package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/c4ts0up/my-stocks/backend/models"
	"gorm.io/gorm"
	"io"
	"net/http"
)

type IFetcherApi interface {
	FetchStockData(url string) ([]models.StockRating, string, error)
	SaveStockData([]models.StockRatingRaw) error
}

// BasicStockFetcher implements a basic fetch
type BasicStockFetcher struct {
	DB *gorm.DB
}

// FetchStockData pulls data from an API and converts it
func (s *BasicStockFetcher) FetchStockData(url string) ([]models.StockRating, string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch data from %v", url)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, "", errors.New("received invalid response from API")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", errors.New("failed to read response body")
	}

	// obtains the stock query response
	var apiResponses models.StockQueryResponse
	if err := json.Unmarshal(body, &apiResponses); err != nil {
		return nil, "", errors.New("failed to parse JSON response")
	}

	// iterates through the stocks
	var stockRatings []models.StockRating
	for _, rawStock := range apiResponses.Stocks {
		stock, err := convertApiResponse(rawStock)
		if err != nil {
			return nil, "", fmt.Errorf("failed to parse stock data, got %v", err)
		}
		stockRatings = append(stockRatings, stock)
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, "", err
	}

	return stockRatings, apiResponses.NextPage, nil
}

// SaveStockData saves stock data to the database (stub implementation). Supposes there are no conflicts in the API
func (s *BasicStockFetcher) SaveStockData(stockList []models.StockRating) error {

	for _, stock := range stockList {
		if err := s.DB.Create(&stock).Error; err != nil {
			return err
		}
	}

	return nil
}
