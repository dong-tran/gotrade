package custom

type ReadSymbols interface {
	Read() []string
}

func NewSymbolsReader(source string) ReadSymbols {
	if source == "env" {
		return NewEnvSymbols()
	} else {
		return NewFileSymbols()
	}
}
