package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port                          string        `mapstructure:"SERVER_PORT"`
	LogLevel                      string        `mapstructure:"LOG_LEVEL"`
	ExpirationWindowInSeconds     time.Duration `mapstructure:"EXPIRATION_WINDOW_IN_SECONDS"`
	SyncFrequencyForPopularTokens time.Duration `mapstructure:"SYNC_FREQUENCY_FOR_POPULAR_TOKENS"`
	SyncPriceErrorCountLimit      int           `mapstructure:"SYNC_PRICE_ERROR_COUNT_LIMIT"`
	PageSizeLimit                 int32         `mapstructure:"PAGE_SIZE_LIMIT"`
}

const (
	expirationWindowInSeconds     = 10
	syncFrequencyForPopularTokens = 60
	syncPriceErrorCountLimit      = 5
	pageSizeLimit                 = 1000
)

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()

	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("LOG_LEVEL", "DEBUG")
	viper.SetDefault("EXPIRATION_WINDOW_IN_SECONDS", expirationWindowInSeconds)
	viper.SetDefault("SYNC_FREQUENCY_FOR_POPULAR_TOKENS", syncFrequencyForPopularTokens)
	viper.SetDefault("SYNC_PRICE_ERROR_COUNT_LIMIT", syncPriceErrorCountLimit)
	viper.SetDefault("PAGE_SIZE_LIMIT", pageSizeLimit)

	var cfg Config
	err := viper.Unmarshal(&cfg)

	return &cfg, err
}
