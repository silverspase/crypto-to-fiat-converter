package memory

import (
	"time"

	"crypto-to-fiat-converter/intenral/config"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"crypto-to-fiat-converter/intenral/service/cache"
	memoryMap "crypto-to-fiat-converter/intenral/service/cache/memory/map"
	price "crypto-to-fiat-converter/intenral/service/price_provider"
)

const (
	noNextPageToken int32 = -1
)

type service struct {
	cfg                 *config.Config
	priceProvider       price.Provider
	priceMap            memoryMap.Provider
	tokenList           []price.TokenListItem
	currencyList        []string
	syncPriceErrorCount int
	metrics             metrics // nolint: unused
}

type metrics struct { // nolint: unused
	tokens     map[string]int
	currencies map[string]int
}

func New(cfg *config.Config, priceProvider price.Provider) (cache.Provider, error) {
	s := &service{
		cfg:           cfg,
		priceProvider: priceProvider,
		priceMap:      memoryMap.New(cfg),
	}

	err := s.initTokenList()
	if err != nil {
		return nil, err
	}

	// start background
	go s.syncPrices()

	return s, nil
}

func (s *service) syncPrices() {
	var start time.Time
	for {
		start = time.Now()
		popularTokens, popularCurrencies := s.getMostFrequentTokensAndCurrencies()
		res, err := s.priceProvider.GetMostFrequentRates(popularTokens, popularCurrencies)
		if err != nil {
			zap.L().Error("failed to update rates", zap.Error(err))
			s.syncPriceErrorCount++
			if s.syncPriceErrorCount > s.cfg.SyncPriceErrorCountLimit {
				zap.L().Error("syncPriceErrorCount exceeded")
				// send alert to dev team, we have an issue there
			}
		}

		for token, currencies := range *res {
			for currency, tokenPrice := range currencies {
				s.priceMap.Upsert(token, currency, tokenPrice)
			}
		}

		zap.L().Debug("finished syncing prices", zap.Int64("elapsed ms", time.Since(start).Milliseconds()))

		time.Sleep(s.cfg.SyncFrequencyForPopularTokens * time.Second)
	}

}

func (s *service) getMostFrequentTokensAndCurrencies() (tokens, currencies []string) {
	// for the first call we can use predefined lists,
	// but after we can collect actual tokens and currencies from metrics and be more user oriented
	// TODO return actual data from metrics
	return []string{"ethereum", "bitcoin"}, []string{"usd", "eur"}
}

func (s *service) GetFiatValue(tokenID, vsCurrency string) (tokenPrice float32, err error) {
	tokenPrice, ok := s.priceMap.Get(tokenID, vsCurrency)
	if !ok {
		zap.L().Info("no cache record, fetching data from price provider")
		return s.fetchTokenPriceAndUpdateStore(tokenID, vsCurrency)
	}

	zap.L().Info("fetching data from cache")

	return tokenPrice, nil
}

func (s *service) GetTokenList(pageToken, pageSize int32) ([]price.TokenListItem, int32, error) {
	start := int(pageToken * pageSize)
	end := start + int(pageSize)

	if len(s.tokenList) <= start {
		return nil, noNextPageToken, status.Errorf(codes.OutOfRange, "page token is too big")
	}

	if len(s.tokenList) <= end {
		return s.tokenList[start:], noNextPageToken, nil
	}

	pageToken++

	return s.tokenList[start:end], pageToken, nil
}

func (s *service) fetchTokenPriceAndUpdateStore(tokenID, vsCurrency string) (tokenPrice float32, err error) {
	tokenPrice, err = s.priceProvider.GetSingleFiatValue(tokenID, vsCurrency)
	if err != nil {
		return 0, err
	}

	// there is no need to wait for upsert completion, so we can send the result to user immediately
	go s.priceMap.Upsert(tokenID, vsCurrency, tokenPrice)

	return tokenPrice, nil
}

func (s *service) initTokenList() (err error) {
	s.tokenList, s.currencyList, err = s.priceProvider.GetTokenAndCurrencyLists()
	return err
}
