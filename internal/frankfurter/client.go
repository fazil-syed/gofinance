package frankfurter

import (
	"encoding/json"
	"fazil-syed/gofinance/internal/config"
	"fmt"
	"net/http"
)

type FrankFurterClient struct {
	cfg *config.FrankFurterConfig
}

func NewFrankFurterClient(cfg *config.FrankFurterConfig) *FrankFurterClient {
	f := &FrankFurterClient{cfg: cfg}
	return f
}

func (f *FrankFurterClient) GetForexAtDate(baseCurrency string, date string) ([]Rate, error) {
	requestUrl := fmt.Sprintf("%s/v2/rates?base=%s&date=%s", f.cfg.BaseURL, baseCurrency, date)

	resp, err := http.Get(requestUrl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var forexRates []Rate

	if err := json.NewDecoder(resp.Body).Decode(&forexRates); err != nil {
		return nil, err
	}

	return forexRates, nil
}
