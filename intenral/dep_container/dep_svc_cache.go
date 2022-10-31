package dep_container

import (
	"crypto-to-fiat-converter/intenral/config"
	"crypto-to-fiat-converter/intenral/service/cache/memory"
	price "crypto-to-fiat-converter/intenral/service/price_provider"
	"github.com/sarulabs/di"
	"go.uber.org/zap"
)

const cacheServiceDefName = "cacheServiceDefName"

// RegisterCacheService registers cache service
func RegisterCacheService(builder *di.Builder) error {
	return builder.Add(di.Def{
		Name: cacheServiceDefName,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get(configDefName).(*config.Config)
			priceProvider := ctn.Get(priceProviderServiceDefName).(price.Provider)

			zap.L().Info("cache service has been started",
				zap.Duration("SyncFrequencyForPopularTokens", cfg.SyncFrequencyForPopularTokens),
				zap.Duration("ExpirationWindowInSeconds", cfg.ExpirationWindowInSeconds),
			)

			// Here we can return different cache services based on Config
			return memory.New(cfg, priceProvider)
		},
	})
}
