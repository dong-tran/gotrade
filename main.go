package main

import (
	"os"
	"strconv"

	"github.com/dong-tran/gotrade/custom"
)

func main() {
	days, err := strconv.Atoi(os.Getenv("LOOKBACK_DAYS"))
	if err != nil {
		days = 730
	}

	fetchSource := os.Getenv("FETCH_SOURCE")
	if len(fetchSource) == 0 {
		fetchSource = "24hmoney"
	}

	symbols := custom.ReadSymbols()
	fetcher := custom.NewFetchData(fetchSource)
	fetcher.FetchData(symbols, days)
	custom.Backtest(days, symbols)
	// Calculate(symbols)
	custom.UploadFile()
}
