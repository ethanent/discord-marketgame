package main

import av "github.com/cmckee-dev/go-alpha-vantage"

func getLivePrice(symbol string) (float64, error) {
	tsvs, err := avClient.StockTimeSeries(av.TimeSeriesDaily, symbol)

	if err != nil {
		return 0, err
	}

	return (*tsvs[len(tsvs)-1]).Close, nil
}
