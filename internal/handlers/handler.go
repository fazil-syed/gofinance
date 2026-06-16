package handlers

import (
	"fazil-syed/gofinance/internal/config"
	finhub "fazil-syed/gofinance/internal/finHub"
	"fazil-syed/gofinance/internal/frankfurter"
)

type Handler struct {
	config            *config.Config
	frankFurterClient *frankfurter.FrankFurterClient
	finhubClient      *finhub.FinHubClient
}

func NewHandler(cfg *config.Config) *Handler {
	finhubClient := finhub.NewFinHubClient(cfg.Finnhub.APIKey, cfg.Finnhub.BaseURL)
	frankfurterClient := frankfurter.NewFrankFurterClient(&cfg.FrankFurter)
	h := &Handler{
		config:            cfg,
		finhubClient:      finhubClient,
		frankFurterClient: frankfurterClient,
	}
	return h
}
