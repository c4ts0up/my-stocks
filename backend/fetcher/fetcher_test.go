package fetcher

import (
	"github.com/c4ts0up/my-stocks/backend/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- TEST CASE 1: INVALID URL ---
func TestFetchStockData_InvalidURL(t *testing.T) {
	fetcher := BasicStockFetcher{}

	_, _, err := fetcher.FetchStockData(":://bad-url")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch data")
}

// --- TEST CASE 2: DOWNSTREAM ERROR ---
func TestFetchStockData_WrongResponse(t *testing.T) {
	// Simulates a server returning 500
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	fetcher := BasicStockFetcher{}
	_, _, err := fetcher.FetchStockData(server.URL)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "received invalid response from API")
}

// --- TEST CASE 3: FAILED TO READ RESPONSE BODY ---
func TestFetchStockData_ReadBodyError(t *testing.T) {
	// Simulate a server that closes the connection immediately after sending headers
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		conn, _, _ := w.(http.Hijacker).Hijack()
		_ = conn.Close()
	}))
	defer server.Close()

	fetcher := BasicStockFetcher{}

	_, _, err := fetcher.FetchStockData(server.URL)
	if err == nil || err.Error() != "failed to read response body" {
		t.Fatalf("expected 'failed to read response body', got %v", err)
	}
}

// --- TEST CASE 4: FAILED TO PARSE JSON RESPONSE ---
func TestFetchStockData_ParseJsonError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Send invalid JSON
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{invalid json}`))
	}))
	defer server.Close()

	fetcher := BasicStockFetcher{}

	_, _, err := fetcher.FetchStockData(server.URL)
	if err == nil || err.Error() != "failed to parse JSON response" {
		t.Fatalf("expected 'failed to parse JSON response', got %v", err)
	}
}

// --- TEST CASE 5: OK FETCH ---
func TestFetchStockData_OK(t *testing.T) {
	mockResponse := `{
		"items": [
			{
				"ticker": "BSBR",
				"target_from": "$4.20",
				"target_to": "$4.70",
				"company": "Banco Santander (Brasil)",
				"action": "upgraded by",
				"brokerage": "The Goldman Sachs Group",
				"rating_from": "Sell",
				"rating_to": "Neutral",
				"time": "2025-01-13T00:30:05.813548892Z"
			}
		],
		"next_page": "AZEK"
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	fetcher := BasicStockFetcher{}

	stocks, nextPage, err := fetcher.FetchStockData(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(stocks) != 1 {
		t.Fatalf("expected 1 stock, got %d", len(stocks))
	}

	if stocks[0].Ticker != "BSBR" {
		t.Errorf("expected ticker BSBR, got %s", stocks[0].Ticker)
	}

	if nextPage != "AZEK" {
		t.Errorf("expected next_page AZEK, got %s", nextPage)
	}
}

// --- TEST CASE 6: SaveStockData works ---
func TestSaveStockData(t *testing.T) {
	// Use an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	_ = db.AutoMigrate(&models.StockRating{})

	fetcher := BasicStockFetcher{DB: db}

	stockList := []models.StockRating{
		{Ticker: "AAPL", Company: "Apple Inc."},
		{Ticker: "GOOGL", Company: "Alphabet Inc."},
	}

	err = fetcher.SaveStockData(stockList)
	if err != nil {
		t.Fatalf("unexpected error saving data: %v", err)
	}

	var count int64
	db.Model(&models.StockRating{}).Count(&count)
	if count != 2 {
		t.Errorf("expected 2 rows in database, got %d", count)
	}
}
