package config

type Config struct {
	Finnhub FinnhubConfig `mapstructure:"Finnhub"`
	Auth    AuthConfig    `mapstructure:"AuthConfig"`
}

type FinnhubConfig struct {
	APIKey  string `mapstructure:"ApiKey"`
	BaseURL string `mapstructure:"BaseURL"`
}

type AuthConfig struct {
	UserName string `mapstructure:"UserName"`
	PassWord string `mapstructure:"PassWord"`
}
