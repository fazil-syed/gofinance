package main

import (
	"fazil-syed/gofinance/internal/config"
	finhub "fazil-syed/gofinance/internal/finHub"
	"fmt"
	"log"
	"time"

	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
)

func main() {

	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	finhubClient := finhub.NewFinHubClient(cfg.Finnhub.APIKey, cfg.Finnhub.BaseURL)

	symbol, err := finhubClient.SearchSymbol("Apple")
	if err != nil {
		log.Fatalf("Error %v", err)
	}

	fmt.Printf("Symbol : %s\n", symbol)
	date := "2026-06-12T11:20:00Z"

	dateParsed, err := time.Parse(time.RFC3339, date)

	if err != nil {
		log.Fatalf("Error %v", err)
	}
	startOfDay := dateParsed.Truncate(24 * time.Hour)

	endOfDay := startOfDay.AddDate(0, 0, 1).Add(-time.Second)

	params := &chart.Params{Symbol: symbol, Interval: datetime.OneDay, Start: datetime.FromUnix(int(startOfDay.Unix())), End: datetime.FromUnix(int(endOfDay.Unix()))}
	iter := chart.Get(params)
	iter.Next()
	price := iter.Bar().Close

	if err := iter.Err(); err != nil {
		log.Fatalf("Error %v", err)
	}

	// price, err := finhubClient.GetStockPrice(symbol, date)

	// if err != nil {
	// 	log.Fatalf("Error %v", err)
	// }

	floatPrice, _ := price.Float64()

	fmt.Printf("Price: %2.f\n", floatPrice)
}
