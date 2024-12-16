package custom

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// SymbolData represents the structure of the JSON response
type symbolDataCafef struct {
	Data struct {
		Value struct {
			DataInfor []struct {
				Time   float64 `json:"time"`
				Open   float64 `json:"open"`
				High   float64 `json:"high"`
				Low    float64 `json:"low"`
				Close  float64 `json:"close"`
				Volume float64 `json:"volume"`
			} `json:"dataInfor"`
		} `json:"value"`
	} `json:"data"`
}

type Cafef struct{}

func (f *Cafef) FetchData(symbols []string, days int) {
	var csvPath = "/tmp/csv/"

	err := os.MkdirAll(csvPath, 0o700)
	if err != nil {
		log.Fatalf("unable to make the output directory: %v", err)
	}

	err = os.RemoveAll(csvPath + "*")
	if err != nil {
		log.Fatalf("Error when clean report output folder: %v", err)
	}

	// Base API URL
	baseURL := "https://msh-devappdata.cafef.vn/rest-api/api/v1/TradingViewsData?symbol=%s&type=D1"

	headers := map[string]string{
		"accept":             "*/*",
		"accept-language":    "en-US,en;q=0.9,vi;q=0.8",
		"cache-control":      "no-cache",
		"origin":             "https://liveboard.cafef.vn",
		"pragma":             "no-cache",
		"priority":           "u=1, i",
		"referer":            "https://liveboard.cafef.vn/",
		"sec-ch-ua":          `"Microsoft Edge";v="131", "Chromium";v="131", "Not_A Brand";v="24"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"macOS"`,
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-site",
		"user-agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0",
	}

	// Loop through the symbols
	for _, symbol := range symbols {
		fmt.Printf("Fetching data for symbol: %s\n", symbol)

		req, err := http.NewRequest("GET", fmt.Sprintf(baseURL, symbol), nil)
		if err != nil {
			fmt.Printf("Error creating request for %s: %v\n", symbol, err)
			continue
		}

		// Set headers
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		/// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error fetching data for %s: %v\n", symbol, err)
			continue
		}
		defer resp.Body.Close()

		// Parse the JSON response
		var symbolData symbolDataCafef
		if err := json.NewDecoder(resp.Body).Decode(&symbolData); err != nil {
			fmt.Printf("Error decoding JSON for %s: %v\n", symbol, err)
			continue
		}

		// Create a CSV file for the symbol
		file, err := os.Create(fmt.Sprintf("%s%s.csv", csvPath, symbol))
		if err != nil {
			fmt.Printf("Error creating CSV file for %s: %v\n", symbol, err)
			continue
		}
		defer file.Close()

		// Write CSV header
		writer := csv.NewWriter(file)
		defer writer.Flush()
		writer.Write([]string{"Date", "Open", "High", "Low", "Close", "Volume"})

		// Write data rows to the CSV
		for i := len(symbolData.Data.Value.DataInfor) - 1; i >= 0; i-- {
			record := symbolData.Data.Value.DataInfor[i]
			date := time.Unix(int64(record.Time), 0).Format("2006-01-02")
			row := []string{
				date,
				fmt.Sprintf("%.2f", record.Open),
				fmt.Sprintf("%.2f", record.High),
				fmt.Sprintf("%.2f", record.Low),
				fmt.Sprintf("%.2f", record.Close),
				// fmt.Sprintf("%.2f", record.Close), // Adj Close is same as Close in this case
				fmt.Sprintf("%.2f", record.Volume),
			}
			writer.Write(row)
		}

		fmt.Printf("Data for %s saved to %s.csv\n", symbol, symbol)
	}
}
