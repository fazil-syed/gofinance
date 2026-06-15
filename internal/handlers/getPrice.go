package handlers

import (
	"encoding/json"
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}

	// res := fmt.Sprintf("Symbol : %s\nCurrency : %s\nExchange Name: %s\nPrice: %2.f\n", symbol, currency, exchangeName, floatPrice)

}
