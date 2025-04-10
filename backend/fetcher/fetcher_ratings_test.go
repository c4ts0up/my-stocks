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

const mockTocken = "mock-token-123"

// --- TEST CASE 1: INVALID URL ---
func TestFetchStockData_InvalidURL(t *testing.T) {
	fetcher := BasicStockRatingsFetcher{BearerToken: mockTocken}

	_, _, _, err := fetcher.FetchStockRatings(":://bad-url")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create request")
}

// --- TEST CASE 2: DOWNSTREAM ERROR ---
func TestFetchStockData_WrongResponse(t *testing.T) {
	// Simulates a server returning 500
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	fetcher := BasicStockRatingsFetcher{BearerToken: mockTocken}
	_, _, _, err := fetcher.FetchStockRatings(server.URL)

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

	fetcher := BasicStockRatingsFetcher{BearerToken: mockTocken}

	_, _, _, err := fetcher.FetchStockRatings(server.URL)
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

	fetcher := BasicStockRatingsFetcher{BearerToken: mockTocken}

	_, _, _, err := fetcher.FetchStockRatings(server.URL)
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

	fetcher := BasicStockRatingsFetcher{BearerToken: mockTocken}

	stocks, nextPage, companyNames, err := fetcher.FetchStockRatings(server.URL)
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

	name, exists := companyNames["BSBR"]
	if exists && name != "Banco Santander (Brasil)" {
		t.Errorf("expected name %s, got %s", "Banco Santander (Brasil)", name)
	}
}

// --- TEST CASE 6: SaveStockRatings works ---
func TestSaveStockData_Ok(t *testing.T) {
	// Use an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	_ = db.AutoMigrate(&models.StockRating{})

	fetcher := BasicStockRatingsFetcher{DB: db, BearerToken: mockTocken}

	stockList := []models.StockRating{
		{Ticker: "AAPL"},
		{Ticker: "GOOGL"},
	}

	err = fetcher.SaveStockRatings(stockList)
	if err != nil {
		t.Fatalf("unexpected error saving data: %v", err)
	}

	var count int64
	db.Model(&models.StockRating{}).Count(&count)
	if count != 2 {
		t.Errorf("expected 2 rows in database, got %d", count)
	}
}

// --- TEST CASE 7: Fetcher without Bearer Token ---
func TestFetchStockData_MissingBearerToken(t *testing.T) {
	// Simulates a server returning 500
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	fetcher := BasicStockRatingsFetcher{}
	_, _, _, err := fetcher.FetchStockRatings(server.URL)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no bearer token provided")
}

// --- TEST CASE 8: FetchAllRatings with multiple pages ---
func TestFetchAll_MultiplePages(t *testing.T) {
	responses := []string{
		`{"items": [{
            "ticker": "BSBR",
            "target_from": "$4.20",
            "target_to": "$4.70",
            "company": "Banco Santander (Brasil)",
            "action": "upgraded by",
            "brokerage": "The Goldman Sachs Group",
            "rating_from": "Sell",
            "rating_to": "Neutral",
            "time": "2025-01-13T00:30:05.813548892Z"
        }], "next_page": "VYGR"}`,
		`{"items": [{
            "ticker": "VYGR",
            "target_from": "$11.00",
            "target_to": "$9.00",
            "company": "Voyager Therapeutics",
            "action": "reiterated by",
            "brokerage": "Wedbush",
            "rating_from": "Outperform",
            "rating_to": "Outperform",
            "time": "2025-01-14T00:30:05.813548892Z"
        }], "next_page": ""}`,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("next_page")
		if page == "VYGR" {
			_, _ = w.Write([]byte(responses[1]))
		} else {
			_, _ = w.Write([]byte(responses[0]))
		}
	}))
	defer server.Close()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	_ = db.AutoMigrate(&models.StockRating{})

	fetcher := BasicStockRatingsFetcher{DB: db, BearerToken: mockTocken}

	var tickers []string
	tickers, err = fetcher.FetchAllRatings(server.URL)
	assert.NoError(t, err)

	var stocks []models.StockRating
	db.Find(&stocks)
	assert.Len(t, stocks, 2)
	assert.Equal(t, "BSBR", stocks[0].Ticker)
	assert.Equal(t, "VYGR", stocks[1].Ticker)
	assert.Equal(t, stocks[0].Ticker, tickers[0])
	assert.Equal(t, stocks[1].Ticker, tickers[1])
}

// --- TEST CASE 9: FetchAllRatings fails on bad page ---
func TestFetchAll_FailsOnBadPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	_ = db.AutoMigrate(&models.StockRating{})

	fetcher := BasicStockRatingsFetcher{DB: db, BearerToken: mockTocken}

	_, err = fetcher.FetchAllRatings(server.URL)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "received invalid response from API")
}

// --- TEST CASE 10: UpdateCompanyNames updates company names
func TestUpdateCompanyNames_Ok(t *testing.T) {
	// Use an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	_ = db.AutoMigrate(&models.Stock{})

	fetcher := BasicStockRatingsFetcher{DB: db, BearerToken: mockTocken}

	companyNames := map[string]string{
		"BSBR": "Banco Santander (Brasil)",
		"VYGR": "",
		"AAPL": "Apple Inc.",
	}

	// This value should not be updated
	db.Where("ticker = VYGR").Create(models.Stock{
		Ticker:  "VYGR",
		Company: "Voyager Therapeutics",
	})

	err = fetcher.UpdateCompanyNames(companyNames)
	assert.NoError(t, err)

	// BSBR is updated with its proper name
	var bsbrStock models.Stock
	db.Where("ticker = ?", "BSBR").First(&bsbrStock)
	assert.Equal(t, companyNames["BSBR"], bsbrStock.Company)

	// AAPL is created and updated
	var aaplStock models.Stock
	db.Where("ticker = ?", "AAPL").First(&aaplStock)
	assert.Equal(t, "Apple Inc.", aaplStock.Company)

	// VYGR is not updated
	var vygrStock models.Stock
	db.Where("ticker = ?", "VYGR").First(&vygrStock)
	assert.Equal(t, "Voyager Therapeutics", vygrStock.Company)
}
