package fetcher

// TODO: fetcher for stock ratings
// TODO: fetcher for stock info
// TODO: full fetcher interface

type StockFetcher struct {
	RatingsFetcher IStockRatingsFetcher
	InfoFetcher    IStockInfoFetcher
}

type IStockRatingsFetcher interface {
	FetchAllRatings(url string) ([]string, error)
}

type IStockInfoFetcher interface {
	FetchAllInfo(tickers []string, url string) error
}

func (f *StockFetcher) FetchAll(ratingsUrl string, infoUrl string) error {
	// fetches rating
	tickers, err := f.RatingsFetcher.FetchAllRatings(ratingsUrl)
	if err != nil {
		return err
	}
	// fetches info
	err = f.InfoFetcher.FetchAllInfo(tickers, infoUrl)
	if err != nil {
		return err
	}

	return nil
}
