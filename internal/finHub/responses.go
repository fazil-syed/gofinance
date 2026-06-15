package finhub

type Symbol struct {
	Description string `json:"description"`
	Symbol      string `json:"symbol"`
}

type SymbolSearchResult struct {
	Count  int      `json:"count"`
	Result []Symbol `json:"result"`
}

type CandlesResult struct {
	Status string    `json:"s"`
	Prices []float32 `json:"c"`
}
