package price

type Provider interface {
	GetSingleFiatValue(tokenName, resultFiat string) (float32, error)
	GetMostFrequentRates(ids, vsCurrencies []string) (*map[string]map[string]float32, error)
	GetTokenAndCurrencyLists() ([]TokenListItem, []string, error)
}
