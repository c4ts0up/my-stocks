package presenter

// StockBase shows the basic information regarding a stock
type StockBase struct {
	Ticker         string  `json:"ticker"`
	CompanyName    string  `json:"company_name"`
	LastPrice      float64 `json:"last_price"`
	Recommendation string  `json:"recommendation"`
}

// StockRating represents the information related to a stock rating
type StockRating struct {
	TargetFrom float64 `json:"target_from"`
	TargetTo   float64 `json:"target_to"`
	Action     string  `json:"action"`
	Brokerage  string  `json:"brokerage"`
	RatingFrom string  `json:"rating_from"`
	RatingTo   string  `json:"rating_to"`
	Time       string  `json:"time"`
}

// StockDetail represents the whole information of a stock and its details
type StockDetail struct {
	StockBase    StockBase     `json:"stock_base"`
	StockRatings []StockRating `json:"stock_ratings"`
}

// StockList gives a base list of all stocks
type StockList []StockBase
