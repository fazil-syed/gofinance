package handlers

import (
	"fazil-syed/gofinance/internal/config"
	finhub "fazil-syed/gofinance/internal/finHub"
	"fazil-syed/gofinance/internal/frankfurter"
	rediscache "fazil-syed/gofinance/internal/redisCache"
)

type Handler struct {
	config            *config.Config
	frankFurterClient *frankfurter.FrankFurterClient
	finhubClient      *finhub.FinHubClient
	redisClient       *rediscache.RedisClient
}

func NewHandler(cfg *config.Config) *Handler {
	finhubClient := finhub.NewFinHubClient(cfg.Finnhub.APIKey, cfg.Finnhub.BaseURL)
	frankfurterClient := frankfurter.NewFrankFurterClient(&cfg.FrankFurter)
	redisClient := rediscache.NewRedisClient(&cfg.Redis)
	h := &Handler{
		config:            cfg,
		finhubClient:      finhubClient,
		frankFurterClient: frankfurterClient,
		redisClient:       redisClient,
	}
	return h
}
