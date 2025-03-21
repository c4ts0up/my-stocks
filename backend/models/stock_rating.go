package models

import "time"
import "gorm.io/gorm"

// StockRating represents a single stock rating update
type StockRating struct {
	gorm.Model               // Includes the ID uint as entry key
	Ticker         string    `json:"ticker"`
	TargetFrom     float64   `json:"target_from"` // TODO: can use uint64 to avoid floating error accumulation
	TargetTo       float64   `json:"target_to"`
	Company        string    `json:"company"`
	Action         string    `json:"action"`
	Brokerage      string    `json:"brokerage"`
	RatingFrom     string    `json:"rating_from"`
	RatingTo       string    `json:"rating_to"`
	Time           time.Time `json:"time"`
	Recommendation string    `json:"recommendation"`
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
