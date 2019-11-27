package main

import (
	"time"

	av "github.com/cmckee-dev/go-alpha-vantage"
)

type tsvl []*av.TimeSeriesValue

var livePxCache map[string]tsvl = map[string]tsvl{}
var lastCachePurge time.Time = time.Now()

func autoPurgeLivePxCache() {
	for {
		time.Sleep(time.Minute * 3)

		livePxCache = map[string]tsvl{}
	}
}

func getLivePrice(symbol string) (float64, error) {
	cachedTsv, ok := livePxCache[symbol]

	if ok {
		return (*cachedTsv[len(cachedTsv)-1]).Close, nil
	}

	tsvs, err := avClient.StockTimeSeries(av.TimeSeriesDaily, symbol)

	if err != nil {
		return 0, err
	}

	lastPx := (*tsvs[len(tsvs)-1]).Close

	livePxCache[symbol] = tsvs

	return lastPx, nil
}

func getDayChange(symbol string) (float64, error) {
	useTsv, ok := livePxCache[symbol]

	if !ok {
		var err error

		useTsv, err = avClient.StockTimeSeries(av.TimeSeriesDaily, symbol)

		if err != nil {
			return -1, err
		}
	}

	var changePercent float64 = ((useTsv[len(useTsv)-1].Close / useTsv[len(useTsv)-2].Close) - 1) * 100

	return changePercent, nil
}
