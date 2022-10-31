package coingecko

import (
	"net/http"
	"time"

	price "crypto-to-fiat-converter/intenral/service/price_provider"
	coingecko "github.com/superoo7/go-gecko/v3"
)

type service struct {
	client *coingecko.Client
}

func New() (price.Provider, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	client := coingecko.NewClient(httpClient)
	_, err := client.Ping()
	if err != nil {
		return nil, err
	}

	return &service{
		client: client,
	}, nil
}

func (s *service) GetSingleFiatValue(tokenID, vsCurrency string) (float32, error) {
	res, err := s.client.SimpleSinglePrice(tokenID, vsCurrency)
	if err != nil {
		return 0, err
	}

	return res.MarketPrice, nil
}

func (s *service) GetFiatValues(ids, vsCurrencies []string) (*map[string]map[string]float32, error) {
	return s.client.SimplePrice(ids, vsCurrencies)
}

func (s *service) GetMostFrequentRates(ids, vsCurrencies []string) (*map[string]map[string]float32, error) {
	return s.client.SimplePrice(ids, vsCurrencies)
}

func (s *service) GetTokenAndCurrencyLists() (tokenList []price.TokenListItem, currencyList []string, err error) {
	tokenListRaw, err := s.client.CoinsList()
	if err != nil {
		return nil, nil, err
	}

	currencyListRaw, err := s.client.SimpleSupportedVSCurrencies()
	if err != nil {
		return nil, nil, err
	}

	for _, val := range *tokenListRaw {
		tokenList = append(tokenList, price.TokenListItem{
			ID:     val.ID,
			Name:   val.Name,
			Symbol: val.Symbol,
		})
	}

	for _, val := range *currencyListRaw {
		currencyList = append(currencyList, val)
	}
	return tokenList, currencyList, nil
}
