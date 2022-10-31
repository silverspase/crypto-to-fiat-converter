package dep_container

import (
	"crypto-to-fiat-converter/intenral/service/cache"
	"crypto-to-fiat-converter/intenral/service/currency_converter"
	"github.com/sarulabs/di"
)

const currencyConverterServiceDefName = "currencyConverterServiceDefName"

// RegisterCurrencyConverterService registers currency converter service
func RegisterCurrencyConverterService(builder *di.Builder) error {
	return builder.Add(di.Def{
		Name: currencyConverterServiceDefName,
		Build: func(ctn di.Container) (interface{}, error) {
			cacheProvider := ctn.Get(cacheServiceDefName).(cache.Provider)

			return currency_converter.New(cacheProvider), nil
		},
	})
}
