package custom

type FetchSource interface {
	FetchData(symbols []string, days int)
}

func NewFetchData(t string) FetchSource {
	if t == "cafef" {
		return &Cafef{}
	} else {
		return &Money24h{}
	}
}
