package fetcher

import "github.com/c4ts0up/my-stocks/backend/models"

type IFetcherApi interface {
	FetchStockData(url string) ([]models.StockRating, error)
	SaveStockData([]models.StockRating) error
}

// BasicStockFetcher implements a basic fetch
type BasicStockFetcher struct{}

func (s *BasicStockFetcher) FetchStockData(url string) ([]models.StockRating, error) {
	return []models.StockRating{}, nil
}

func (s *BasicStockFetcher) SaveStockData([]models.StockRating) error {
	return nil
}
