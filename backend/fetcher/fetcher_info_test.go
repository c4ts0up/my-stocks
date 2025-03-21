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

const mockToken = "mock-token-123"

// --- TEST CASE 1: INVALID URL ---
func TestFetchStockInfo_InvalidURL(t *testing.T) {
	fetcher := BasicStockInfoFetcher{BearerToken: mockToken}

	_, err := fetcher.FetchStockInfo("AAPL", ":://bad-url")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid base URL:")
}

// --- TEST CASE 2: DOWNSTREAM ERROR ---
func TestFetchStockInfo_WrongResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	fetcher := BasicStockInfoFetcher{BearerToken: mockToken}
	_, err := fetcher.FetchStockInfo("AAPL", server.URL)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code")
}

// --- TEST CASE 3: FAILED TO READ RESPONSE BODY ---
func TestFetchStockInfo_ReadBodyError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		conn, _, _ := w.(http.Hijacker).Hijack()
		_ = conn.Close()
	}))
	defer server.Close()

	fetcher := BasicStockInfoFetcher{BearerToken: mockToken}

	_, err := fetcher.FetchStockInfo("AAPL", server.URL)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read response body")
}

// --- TEST CASE 4: FAILED TO PARSE JSON RESPONSE ---
func TestFetchStockInfo_ParseJsonError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{invalid json}`))
	}))
	defer server.Close()

	fetcher := BasicStockInfoFetcher{BearerToken: mockToken}

	_, err := fetcher.FetchStockInfo("AAPL", server.URL)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse JSON")
}

// --- TEST CASE 5: OK FETCH ---
func TestFetchStockInfo_OK(t *testing.T) {
	mockResponse := `[{"ticker":"AAPL","open":211.51,"lastClose":214.1,"lastPrice":214.65,"percentage":0.26,"currency":"USD","companyName":"Apple Inc."}]`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("ticker")
		if query == "AAPL" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(mockResponse))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	fetcher := BasicStockInfoFetcher{BearerToken: mockToken}

	stock, err := fetcher.FetchStockInfo("AAPL", server.URL+"?ticker")
	assert.NoError(t, err)
	assert.Equal(t, "AAPL", stock.Ticker)
	assert.Equal(t, "Apple Inc.", stock.Company)
	assert.Equal(t, 214.65, stock.LastPrice)
}

// --- TEST CASE 6: SaveStockInfo works ---
func TestSaveStockInfo(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	_ = db.AutoMigrate(&models.Stock{})

	fetcher := BasicStockInfoFetcher{DB: db, BearerToken: mockToken}

	stock := models.Stock{Ticker: "AAPL", Company: "Apple Inc."}
	err = fetcher.SaveStockInfo(stock)
	assert.NoError(t, err)

	var count int64
	db.Model(&models.Stock{}).Count(&count)
	assert.Equal(t, int64(1), count)
}

// --- TEST CASE 7: FetchAllInfo with multiple tickers ---
func TestFetchAllInfo_MultipleTickers(t *testing.T) {
	responses := map[string]string{
		"AAPL":  `[{"ticker":"AAPL","companyName":"Apple Inc."}]`,
		"GOOGL": `[{"ticker":"GOOGL","companyName":"Alphabet Inc."}]`,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("ticker")
		if response, ok := responses[query]; ok {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(response))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	_ = db.AutoMigrate(&models.Stock{})

	fetcher := BasicStockInfoFetcher{DB: db, BearerToken: mockToken}
	tickers := []string{"AAPL", "GOOGL"}
	err = fetcher.FetchAllInfo(tickers, server.URL+"?ticker")
	assert.NoError(t, err)

	var stocks []models.Stock
	db.Find(&stocks)
	assert.Len(t, stocks, 2)
	assert.Equal(t, "AAPL", stocks[0].Ticker)
	assert.Equal(t, "GOOGL", stocks[1].Ticker)
}
