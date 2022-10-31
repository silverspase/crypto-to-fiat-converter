package currency_converter

import (
	"crypto-to-fiat-converter/intenral/service/cache"
	"crypto-to-fiat-converter/intenral/service/currency_converter/types"
	price "crypto-to-fiat-converter/intenral/service/price_provider"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	cacheProvider cache.Provider
}

func New(cacheProvider cache.Provider) Provider {
	return &service{
		cacheProvider: cacheProvider,
	}
}

func (s *service) Convert(request []types.ConvertRequest) (response []types.ConvertResponse, err error) {
	var singlePrice float32
	for _, item := range request {
		singlePrice, err = s.cacheProvider.GetFiatValue(item.FromToken, item.ToCurrency)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "failed to fetch fiat value")
		}
		response = append(response, types.ConvertResponse{
			FromToken:   item.FromToken,
			ToCurrency:  item.ToCurrency,
			TokenAmount: item.Amount,
			TotalPrice:  singlePrice * item.Amount,
		})
	}

	return response, nil
}

func (s *service) GetTokenList(pageToken, pageSize int32) ([]price.TokenListItem, int32, error) {
	return s.cacheProvider.GetTokenList(pageToken, pageSize)
}
