package currency_converter

import (
	"crypto-to-fiat-converter/intenral/service/currency_converter/types"
	price "crypto-to-fiat-converter/intenral/service/price_provider"
)

type Provider interface {
	Convert(request []types.ConvertRequest) (response []types.ConvertResponse, err error)
	GetTokenList(pageToken, pageSize int32) ([]price.TokenListItem, int32, error)
}
