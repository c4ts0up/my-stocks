package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/c4ts0up/my-stocks/backend/models"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
)

type BasicStockRatingsFetcher struct {
	DB          *gorm.DB
	BearerToken string
}

// FetchStockRatings pulls stock ratings from the given API and converts them to StockRating models
func (s *BasicStockRatingsFetcher) FetchStockRatings(url string) ([]models.StockRating, string, map[string]string, error) {
	log.Printf("Fetching stock data from %s\n", url)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to create request: %w", err)
	}
	if s.BearerToken == "" {
		return nil, "", nil, fmt.Errorf("no bearer token provided")
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.BearerToken))

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to fetch data from %v: %w", url, err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, "", nil, fmt.Errorf("received invalid response from API: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", nil, errors.New("failed to read response body")
	}

	// obtains the stock query response
	var apiResponses models.StockQueryResponse
	if err := json.Unmarshal(body, &apiResponses); err != nil {
		return nil, "", nil, errors.New("failed to parse JSON response")
	}

	// iterates through the stocks
	var stockRatings []models.StockRating
	for _, rawStock := range apiResponses.Stocks {
		stock, err := convertStockRatingsApiResponse(rawStock)
		if err != nil {
			return nil, "", nil, fmt.Errorf("failed to parse stock data, got %v", err)
		}
		stockRatings = append(stockRatings, stock)
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, "", nil, err
	}

	// Extract all the names
	companyNames := map[string]string{}
	for _, stock := range apiResponses.Stocks {
		companyNames[stock.Ticker] = stock.Company
	}

	return stockRatings, apiResponses.NextPage, companyNames, nil
}

// SaveStockRatings saves StockRating models to the database (stub implementation). Supposes there are no conflicts in the API and Stocks are already created
func (s *BasicStockRatingsFetcher) SaveStockRatings(stockList []models.StockRating) error {
	log.Printf("Saving stock data to database")

	// Upsert each stock record (insert or update)
	for _, stock := range stockList {
		err := s.DB.Where("ticker = ? AND time = ?", stock.Ticker, stock.Time).
			Assign(stock).
			FirstOrCreate(&stock).Error

		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateCompanyNames updates the company name in the database if a non-empty name if found in the API response
func (s *BasicStockRatingsFetcher) UpdateCompanyNames(companyNames map[string]string) error {
	for ticker, name := range companyNames {
		if name != "" {
			result := s.DB.Model(&models.Stock{}).
				Where("ticker = ?", ticker).
				Updates(models.Stock{Company: name})

			if result.Error != nil {
				return result.Error
			}

			// No rows updated. Must create ticker
			if result.RowsAffected == 0 {
				stock := models.Stock{Ticker: ticker, Company: name}
				if err := s.DB.Create(&stock).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// GetStockTickers gets the tickers from the StockRating list obtained after fetching
func (s *BasicStockRatingsFetcher) GetStockTickers(stockRatings []models.StockRating) []string {
	log.Printf("Extracting stock tickers")

	tickers := make([]string, len(stockRatings))
	for i, v := range stockRatings {
		tickers[i] = v.Ticker
	}

	return tickers
}

// FetchAllRatings fetches all the pages in the database and saves them in the ORM
func (s *BasicStockRatingsFetcher) FetchAllRatings(url string) ([]string, error) {
	nextPage := ""
	var tickers []string

	for {

		// pulls
		stockRatings, newSuffix, companyNames, err := s.FetchStockRatings(url + "?next_page=" + nextPage)
		if err != nil {
			log.Printf("Error during stock ratings fetch:  %e", err)
			return []string{}, err
		}

		// saves
		err = s.SaveStockRatings(stockRatings)
		if err != nil {
			log.Printf("Error during stock ratings update in the database: %e", err)
			return []string{}, err
		}

		// updates names in the database
		err = s.UpdateCompanyNames(companyNames)
		if err != nil {
			log.Printf("Error during company names update: %e", err)
		}

		// adds tickers to return
		tickers = append(tickers, s.GetStockTickers(stockRatings)...)

		// checks if there are more pages to fetch
		if newSuffix == "" {
			break
		}
		nextPage = newSuffix
	}

	return tickers, nil
}
