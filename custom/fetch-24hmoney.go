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

type Money24h struct{}

type myTime struct {
	time.Time
}

func (m *myTime) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	tt, err := time.Parse("\"2006-01-02T15:04:05\"", string(data))
	if err == nil {
		*m = myTime{tt}
	} else {
		t2, err2 := time.Parse(`\"`+time.RFC3339+`\"`, string(data))
		if err2 == nil {
			*m = myTime{t2}
		} else {
			t3, err3 := time.Parse("\"2006-01-02T15:04:05-07:00\"", string(data))
			if err3 == nil {
				*m = myTime{t3}
			}
			return err3
		}
	}
	return err
}

type symbolData24h struct {
	Items []struct {
		TradingDate      myTime  `json:"tradingDate"`
		OpenPrice        float64 `json:"openPrice"`
		ClosePrice       float64 `json:"closePrice"`
		HighestPrice     float64 `json:"highestPrice"`
		LowestPrice      float64 `json:"lowestPrice"`
		TotalMatchVolume float64 `json:"totalMatchVolume"`
	} `json:"items"`
}

func (f *Money24h) FetchData(symbols []string, days int) {
	var csvPath = "/tmp/csv/"

	err := os.MkdirAll(csvPath, 0o700)
	if err != nil {
		log.Fatalf("unable to make the output directory: %v", err)
	}

	err = os.RemoveAll(csvPath + "*")
	if err != nil {
		log.Fatalf("Error when clean report output folder: %v", err)
	}

	toDate := time.Now().Truncate(time.Hour)
	fromDate := toDate.AddDate(-1, 0, -days)
	format := "2006-01-02T15:04:05.999Z"

	// Base API URL
	baseURL := "https://24hmoney.vn/TradingView/GetStockChartData?Code=%s&Frequency=Daily&From=%s&To=%s&language=vi&Type=Stock"

	headers := map[string]string{
		"accept":             "*/*",
		"accept-language":    "en-US,en;q=0.9,vi;q=0.8",
		"cache-control":      "no-cache",
		"pragma":             "no-cache",
		"priority":           "u=1, i",
		"referer":            "https://24hmoney.vn/stock/HPG/technical-analysis",
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

		req, err := http.NewRequest("GET", fmt.Sprintf(baseURL, f.getSymbol(symbol), fromDate.Format(format), toDate.Format(format)), nil)
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
		var symbolData symbolData24h
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
		for i := 0; i < len(symbolData.Items); i++ {
			record := symbolData.Items[i]
			date := record.TradingDate.Format("2006-01-02")
			row := []string{
				date,
				fmt.Sprintf("%.2f", record.OpenPrice),
				fmt.Sprintf("%.2f", record.HighestPrice),
				fmt.Sprintf("%.2f", record.LowestPrice),
				fmt.Sprintf("%.2f", record.ClosePrice),
				// fmt.Sprintf("%.2f", record.Close), // Adj Close is same as Close in this case
				fmt.Sprintf("%.2f", record.TotalMatchVolume),
			}
			writer.Write(row)
		}

		fmt.Printf("Data for %s saved to %s.csv\n", symbol, symbol)
	}
}

func (f *Money24h) getSymbol(s string) string {
	symbols := map[string]string{
		"VGI":     "VTGI",
		"SSB":     "SEAB",
		"LPB":     "LVB",
		"BAF":     "BAFVN",
		"YEG":     "YEGCORP",
		"VFS":     "VFSC",
		"DSE":     "DNSE",
		"GEX":     "GELEX",
		"HHV":     "HAMADECO",
		"VEA":     "VEAM",
		"DPG":     "DATPHUONG",
		"FTS":     "FPTS",
		"BFC":     "BDFC",
		"VCI":     "VCSC",
		"HVN":     "VNAIR",
		"IDC":     "IDICO",
		"DXS":     "DXRES",
		"VTZ":     "10780",
		"VJC":     "VIETJET",
		"FRT":     "FPTR",
		"KHG":     "KHAIHOANLAND",
		"KOS":     "KOSY",
		"GVR":     "VNRG",
		"VNINDEX": "VN-INDEX",
		"PSH":     "NSHPETRO",
		"LTG":     "AGPPS",
		"MBS":     "TLSC",
		"DTD":     "PTTD",
		"TCH":     "HHSF",
		"BSR":     "BSRC",
		"SZC":     "SIDC",
		"AGG":     "AGRE",
		"VTP":     "VTPO",
		"ACV":     "ACVN",
		"DBD":     "BIDIP",
		"VHM":     "NHN",
		"PC1":     "PCC1",
		"EVG":     "EVGC",
		"SIP":     "SGVRG",
		"EVF":     "EVNF",
		"POW":     "PVPOWER",
		"VRE":     "VRJSC",
		"SAB":     "SBCO",
		"NVL":     "NOVALAND",
		"VPI":     "VPVICTO",
		"PLX":     "PETRO",
		"PVP":     "PVTRANS",
	}
	symbol, ok := symbols[s]
	if ok {
		return symbol
	}
	return s
}
