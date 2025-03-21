package models

import "time"

// Stock represents an observed stock and its static information
type Stock struct {
	Ticker         string `gorm:"primaryKey"`
	Company        string
	Recommendation string
}

// StockRating represents the most recent stock rating given by some broker
type StockRating struct {
	Ticker     string `gorm:"primaryKey"`
	Brokerage  string `gorm:"primaryKey"`
	TargetFrom float64
	TargetTo   float64
	Action     string
	RatingFrom string
	RatingTo   string
	Time       time.Time
}

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
