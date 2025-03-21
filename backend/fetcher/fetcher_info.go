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

type BasicStockInfoFetcher struct {
	DB          *gorm.DB
	BearerToken string
}

// FetchStockInfo fetches stock data from Algobook Stock API
func (b *BasicStockInfoFetcher) FetchStockInfo(ticker string, url string) (models.Stock, error) {
	finalUrl := fmt.Sprintf("https://stocks.algobook.info/api/v1/stocks?tickers=%s", ticker) // FIXME: unburn

	resp, err := http.Get(finalUrl)
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
		Ticker:  data[0].Ticker,
		Company: data[0].CompanyName,
	}

	return stock, nil
}

// SaveStockInfo saves a Stock model to the database
func (b *BasicStockInfoFetcher) SaveStockInfo(stock models.Stock) error {
	if err := b.DB.Save(&stock).Error; err != nil {
		return fmt.Errorf("failed to save stock data for %s: %w", stock.Ticker, err)
	}
	return nil
}

// FetchAllInfo fetches and saves data for all given tickers
func (b *BasicStockInfoFetcher) FetchAllInfo(tickers []string, url string) error {
	if len(tickers) == 0 {
		return errors.New("no tickers provided")
	}

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
