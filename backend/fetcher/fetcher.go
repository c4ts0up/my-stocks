package fetcher

// TODO: fetcher for stock ratings
// TODO: fetcher for stock info
// TODO: full fetcher interface

type IStockFetcher struct {
	ratingsFetcher *IStockRatingsFetcher
	infoFetcher    *IStockInfoFetcher
}

type IStockRatingsFetcher interface {
	FetchAllRatings(url string) error
}

type IStockInfoFetcher interface {
	FetchAllInfo(tickers []string, url string) error
}
