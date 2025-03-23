package fetcher

import (
	"encoding/json"
	"fmt"
	"github.com/c4ts0up/my-stocks/backend/models"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"net/url"
)

type BasicStockInfoFetcher struct {
	DB          *gorm.DB
	BearerToken string
}

// FetchStockInfo fetches stock data from Algobook Stock API
func (b *BasicStockInfoFetcher) FetchStockInfo(ticker string, baseUrl string) (models.Stock, error) {
	// Parse the base URL
	u, err := url.Parse(baseUrl)
	if err != nil {
		return models.Stock{}, fmt.Errorf("invalid base URL: %w", err)
	}

	// Add the ticker as a query parameter
	q := u.Query()
	q.Set("tickers", ticker)
	u.RawQuery = q.Encode()

	log.Printf("Fetching stock info from %s", u.String())

	resp, err := http.Get(u.String())
	if err != nil {
		return models.Stock{}, fmt.Errorf("failed to fetch data for %s: %w", ticker, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Stock{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Stock{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var data models.StockInfoQueryResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return models.Stock{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(data) == 0 {
		return models.Stock{}, fmt.Errorf("no data returned for ticker %s", ticker)
	}

	// Extract stock details
	stock := models.Stock{
		Ticker:    data[0].Ticker,
		Company:   data[0].CompanyName,
		LastPrice: data[0].LastPrice,
	}

	return stock, nil
}

// SaveStockInfo saves a Stock model to the database
func (b *BasicStockInfoFetcher) SaveStockInfo(stock models.Stock) error {
	result := b.DB.Model(&models.Stock{}).Where("ticker = ?", stock.Ticker).Updates(models.Stock{
		LastPrice: stock.LastPrice,
		Company:   stock.Company,
	})

	// If no rows were affected, it's a new stock, so we insert it
	if result.RowsAffected == 0 {
		if err := b.DB.Create(&stock).Error; err != nil {
			return fmt.Errorf("failed to insert new stock data for %s: %w", stock.Ticker, err)
		}
	}

	return nil
}

// FetchAllInfo fetches and saves data for all given tickers
func (b *BasicStockInfoFetcher) FetchAllInfo(tickers []string, url string) error {

	for _, ticker := range tickers {
		stock, err := b.FetchStockInfo(ticker, url)
		if err != nil {
			return fmt.Errorf("failed to fetch data for ticker %s: %w", ticker, err)
		}

		if err := b.SaveStockInfo(stock); err != nil {
			return fmt.Errorf("failed to save data for ticker %s: %w", ticker, err)
		}
	}

	return nil
}
