package config

type Config struct {
	Finnhub     FinnhubConfig     `mapstructure:"Finnhub"`
	Auth        AuthConfig        `mapstructure:"AuthConfig"`
	Server      ServerConfig      `mapstructure:"ServerConfig"`
	FrankFurter FrankFurterConfig `mapstructure:"FrankFurterConfig"`
	Redis       RedisConfig       `mapstructure:"RedisConfig"`
}

type ServerConfig struct {
	Port int `mapstructure:"Port" validate:"required,gte=0"`
}

type FinnhubConfig struct {
	APIKey  string `mapstructure:"ApiKey" validate:"required"`
	BaseURL string `mapstructure:"BaseURL" validate:"required,url"`
}

type AuthConfig struct {
	UserName string `mapstructure:"UserName" validate:"required"`
	PassWord string `mapstructure:"PassWord" validate:"required"`
}

type FrankFurterConfig struct {
	BaseURL string `mapstructure:"BaseURL" validate:"required,url"`
}

type RedisConfig struct {
	Address  string `mapstructure:"Address" validate:"required"`
	PassWord string `mapstructure:"PassWord" validate:"required"`
}
