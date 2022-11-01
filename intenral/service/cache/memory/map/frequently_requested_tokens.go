package memory_map

import (
	"sync"
	"time"

	"go.uber.org/zap"

	"crypto-to-fiat-converter/intenral/config"
)

type pricesByTokenAndCurrency struct {
	cfg *config.Config
	m   map[string]map[string]tokenPrice // map[token]map[currency]tokenPrice
	mu  sync.RWMutex
}

type tokenPrice struct {
	price     float32
	updatedAt time.Time
}

func New(cfg *config.Config) *pricesByTokenAndCurrency {
	return &pricesByTokenAndCurrency{
		cfg: cfg,
		m:   make(map[string]map[string]tokenPrice),
		mu:  sync.RWMutex{},
	}
}

func (p *pricesByTokenAndCurrency) Upsert(token, currency string, price float32) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.m[token] == nil {
		p.m[token] = make(map[string]tokenPrice)
	}

	p.m[token][currency] = tokenPrice{
		price:     price,
		updatedAt: time.Now(),
	}
}

func (p *pricesByTokenAndCurrency) Get(token, currency string) (float32, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.m[token] == nil {
		return 0, false
	}

	res, ok := p.m[token][currency]
	// if there is an entry for given token&currency,
	// but it was updated more than ExpirationWindowInSeconds ago, remove entry and return false to trigger price_provider
	expirationTime := res.updatedAt.Add(p.cfg.ExpirationWindowInSeconds * time.Second)
	if ok && time.Now().After(expirationTime) {
		zap.L().Info("price is too old, removing it")
		delete(p.m[token], currency)

		return 0, false
	}

	return res.price, ok
}
