package custom

import (
	"os"
	"strings"
)

type envSymbols struct{}

func NewEnvSymbols() ReadSymbols {
	return &envSymbols{}
}

func (s envSymbols) Read() []string {
	source := os.Getenv("SYMBOLS_TARGET")
	return strings.Split(source, ",")
}
