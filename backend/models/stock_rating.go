package models

import "time"

// Stock represents an observed stock and its static information
type Stock struct {
	Ticker         string `gorm:"primaryKey"`
	LastPrice      float64
	Company        string
	Recommendation string
}

// StockRating represents the most recent stock rating given by some broker
type StockRating struct {
	Ticker     string `gorm:"primaryKey"` // FIXME: add foreign key
	Brokerage  string `gorm:"primaryKey"`
	TargetFrom float64
	TargetTo   float64
	Action     string
	RatingFrom string
	RatingTo   string
	Time       time.Time
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///			MY STOCKS DOWNSTREAM API MODELS
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// StockRatingRaw matches the raw stock structure in a response
type StockRatingRaw struct {
	Ticker     string `json:"ticker"`
	TargetFrom string `json:"target_from"`
	TargetTo   string `json:"target_to"`
	Company    string `json:"company"`
	Action     string `json:"action"`
	Brokerage  string `json:"brokerage"`
	RatingFrom string `json:"rating_from"`
	RatingTo   string `json:"rating_to"`
	Time       string `json:"time"`
}

// StockQueryResponse matches the JSON response when querying stocks
type StockQueryResponse struct {
	Stocks   []StockRatingRaw `json:"items"`
	NextPage string           `json:"next_page"`
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///			STOCK INFO API MODELS
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// StockInfoRaw represents the raw API response structure for each stock.
type StockInfoRaw struct {
	Ticker      string  `json:"ticker"`
	Open        float64 `json:"open"`
	LastClose   float64 `json:"lastClose"`
	LastPrice   float64 `json:"lastPrice"`
	Percentage  float64 `json:"percentage"`
	Currency    string  `json:"currency"`
	CompanyName string  `json:"companyName"`
}

// StockInfoQueryResponse represents the API response containing multiple stocks.
type StockInfoQueryResponse []StockInfoRaw
