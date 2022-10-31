package cache

import price "crypto-to-fiat-converter/intenral/service/price_provider"

type Provider interface {
	GetFiatValue(tokenName, resultFiat string) (float32, error)
	GetTokenList(pageToken, pageSize int32) ([]price.TokenListItem, int32, error)
}
