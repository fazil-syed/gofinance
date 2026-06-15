package handlers

import (
	"fazil-syed/gofinance/internal/config"
	finhub "fazil-syed/gofinance/internal/finHub"
)

type Handler struct {
	config *config.Config

	finhubClient *finhub.FinHubClient
}

func NewHandler(cfg *config.Config) *Handler {
	finhubClient := finhub.NewFinHubClient(cfg.Finnhub.APIKey, cfg.Finnhub.BaseURL)
	h := &Handler{
		config:       cfg,
		finhubClient: finhubClient,
	}
	return h
}
