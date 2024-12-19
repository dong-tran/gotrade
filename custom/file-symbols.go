package custom

import (
	"embed"
	"encoding/json"
	"log"
)

//go:embed data/symbols.json
var symbolsFS embed.FS

type fileSymbol struct{}

func NewFileSymbols() ReadSymbols {
	return &fileSymbol{}
}

func (s *fileSymbol) Read() []string {
	var symbols []string

	data, err := symbolsFS.ReadFile("data/symbols.json")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// Unmarshal the JSON data into the 'symbols' slice
	err = json.Unmarshal(data, &symbols)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	return symbols
}
