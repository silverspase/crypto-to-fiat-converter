package dep_container

import (
	"crypto-to-fiat-converter/intenral/service/price_provider/coingecko"
	"github.com/sarulabs/di"
)

const priceProviderServiceDefName = "priceProviderServiceDefName"

// RegisterPriceProviderService registers cache service
func RegisterPriceProviderService(builder *di.Builder) error {
	return builder.Add(di.Def{
		Name: priceProviderServiceDefName,
		Build: func(ctn di.Container) (interface{}, error) {

			// Here we can return different price providers based on Config
			return coingecko.New()
		},
	})
}
