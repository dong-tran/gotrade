package main

import (
	"os"
	"strconv"

	"github.com/dong-tran/gotrade/custom"
)

func main() {
	days, err := strconv.Atoi(os.Getenv("LOOKBACK_DAYS"))
	if err != nil {
		days = 95
	}

	fetchSource := os.Getenv("FETCH_SOURCE")
	if len(fetchSource) == 0 {
		fetchSource = "24hmoney"
	}

	symbolsSource := os.Getenv("SYMBOLS_SOURCE")
	symbolsReader := custom.NewSymbolsReader(symbolsSource)
	symbols := symbolsReader.Read()

	fetcher := custom.NewFetchData(fetchSource)
	fetcher.FetchData(symbols, 730)
	custom.Backtest(days, symbols)
	// // Calculate(symbols)
	custom.UploadFile()
}
