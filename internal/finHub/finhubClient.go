package finhub

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type FinHubClient struct {
	Client http.Client

	Key string

	BaseURL string
}

func NewFinHubClient(key string, baseURL string) *FinHubClient {
	f := &FinHubClient{Key: key, BaseURL: baseURL}
	f.Client = http.Client{
		Timeout: time.Second * 30,
	}
	return f
}

func (f *FinHubClient) Call(method string, url string, response any) error {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return err
	}
	req.Header.Set("X-Finnhub-Token", f.Key)
	resp, err := f.Client.Do(req)

	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Non 200 status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(response)

	if err != nil {
		return err
	}
	return nil

}

func (f *FinHubClient) SearchSymbol(company string) (string, error) {
	searchURL := fmt.Sprintf("%s/search?q=%s", f.BaseURL, url.QueryEscape(company))
	var searchResult SymbolSearchResult

	err := f.Call("GET", searchURL, &searchResult)
	if err != nil {
		return "", err
	}

	if len(searchResult.Result) <= 0 {
		return "", errors.New("Symbol Not Found")
	}

	return searchResult.Result[0].Symbol, nil
}

func (f *FinHubClient) GetStockPrice(symbol string, dateStr string) (float32, error) {

	date, err := time.Parse(time.RFC3339, dateStr)

	if err != nil {
		return 0.0, err
	}

	startOfDay := date.Truncate(24 * time.Hour)

	endOfDay := startOfDay.AddDate(0, 0, 1).Add(-time.Second)

	searchURL := fmt.Sprintf("%s/stock/candle?symbol=%s&resolution=D&from=%d&to=%d", f.BaseURL, symbol, startOfDay.Unix(), endOfDay.Unix())

	var candleResult CandlesResult

	err = f.Call("GET", searchURL, &candleResult)
	if err != nil {
		return 0.0, err
	}

	if len(candleResult.Prices) <= 0 {
		return 0.0, errors.New("Prices not found")
	}

	return candleResult.Prices[0], nil
}
