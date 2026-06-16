package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
)

type StockResponse struct {
	Price    float64 `json:"price"`
	Symbol   string  `json:"symbol"`
	Exchange string  `json:"exchange"`
	Currency string  `json:"currency"`
}

var RedisTTL = 8766 * time.Hour

// GetPriceAtDayHandler godoc
//
//	@Security		BasicAuth
//
//	@Summary		Get price of a stock at a given day
//	@Description	Fetch closing price of a company on any past date
//	@Tags			Stock
//	@Produce		json
//	@Param			company	query		string	false	"Company Name "
//	@Param			date	query		string	false	"Date to fetch price at"
//	@Success		200		{object}	StockResponse
//	@Router			/price [get]
func (h *Handler) GetPriceAtDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	company := r.URL.Query().Get("company")

	if company == "" {
		http.Error(w, "missing company param", http.StatusBadRequest)
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		http.Error(w, "missing date param", http.StatusBadRequest)
		return
	}
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	key := fmt.Sprintf("gofinance:%s:%s", company, dateStr)

	//search in cache
	cached, err := h.redisClient.GetData(r.Context(), key)

	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(cached)
		return
	}

	symbol, err := h.finhubClient.SearchSymbol(company)

	if err != nil {
		http.Error(w, "company symbol not found", http.StatusNotFound)
		return
	}
	startOfDay := parsedDate.Truncate(24 * time.Hour)

	endOfDay := startOfDay.AddDate(0, 0, 1).Add(-time.Second)

	params := &chart.Params{Symbol: symbol, Interval: datetime.OneDay, Start: datetime.FromUnix(int(startOfDay.Unix())), End: datetime.FromUnix(int(endOfDay.Unix()))}
	iter := chart.Get(params)
	if !iter.Next() {
		http.Error(w, "No Data found for given date", http.StatusNotFound)
		return
	}
	price := iter.Bar().Close

	currency := iter.Meta().Currency

	exchangeName := iter.Meta().ExchangeName

	if err := iter.Err(); err != nil {
		http.Error(w, "Error Fetching Price", http.StatusInternalServerError)
		return
	}
	floatPrice, _ := price.Float64()
	response := &StockResponse{
		Price:    floatPrice,
		Symbol:   symbol,
		Exchange: exchangeName,
		Currency: currency,
	}
	// save to cache

	data, err := json.Marshal(response)

	if err == nil {
		h.redisClient.SetData(r.Context(), key, data, RedisTTL)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}

// GetForexRatesAtDate godoc
//
//	@Security		BasicAuth
//
//	@Summary		Get forex rates of currencies against a base currency at a given day
//	@Description	Fetch forex rates of 165 currencies against a base currency at a given day from 1949
//	@Tags			Forex
//	@Produce		json
//	@Param			base_currency	query	string	false	"Base Currency "
//	@Param			date			query	string	false	"Date to fetch forex rates at"
//	@Success		200				{array}	frankfurter.Rate
//	@Router			/forex [get]
func (h *Handler) GetForexRatesAtDate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	baseCurrency := r.URL.Query().Get("base_currency")
	if baseCurrency == "" {
		http.Error(w, "missing base_currency param", http.StatusBadRequest)
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		http.Error(w, "missing date param", http.StatusBadRequest)
		return
	}
	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	key := fmt.Sprintf("gofinance:%s:%s", baseCurrency, dateStr)

	//search in cache
	cached, err := h.redisClient.GetData(r.Context(), key)

	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(cached)
		return
	}

	forexRates, err := h.frankFurterClient.GetForexAtDate(baseCurrency, dateStr)

	if err != nil {
		fmt.Printf("%v", err)
		http.Error(w, "Error fetching forex rates", http.StatusInternalServerError)
		return
	}
	// save to cache

	data, err := json.Marshal(forexRates)

	if err == nil {
		h.redisClient.SetData(r.Context(), key, data, RedisTTL)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(forexRates); err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}
