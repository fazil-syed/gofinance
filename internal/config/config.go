package config

type Config struct {
	Finnhub FinnhubConfig `mapstructure:"Finnhub"`
}

type FinnhubConfig struct {
	APIKey  string `mapstructure:"ApiKey"`
	BaseURL string `mapstructure:"BaseURL"`
}
