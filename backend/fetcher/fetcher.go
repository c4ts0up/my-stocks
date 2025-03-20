package fetcher

import (
	"encoding/json"
	"errors"
	"github.com/c4ts0up/my-stocks/backend/models"
	"io"
	"net/http"
)

type IFetcherApi interface {
	FetchStockData(url string) ([]models.StockRating, string, error)
	SaveStockData([]models.StockRatingRaw) error
}

// BasicStockFetcher implements a basic fetch
type BasicStockFetcher struct{}

// FetchStockData pulls data from an API and converts it
func (s *BasicStockFetcher) FetchStockData(url string) ([]models.StockRating, string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, "", errors.New("failed to fetch data")
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
			return nil, "", err
		}
		stockRatings = append(stockRatings, stock)
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, "", err
	}

	return stockRatings, apiResponses.NextPage, nil
}

// SaveStockData saves stock data to the database (stub implementation)
func (s *BasicStockFetcher) SaveStockData(stockList []models.StockRating) error {
	// TODO: Implement this with GORM later
	return nil
}
