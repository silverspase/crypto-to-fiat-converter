package grpc_service

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"crypto-to-fiat-converter/intenral/config"
	"crypto-to-fiat-converter/intenral/service/currency_converter"
	"crypto-to-fiat-converter/intenral/service/currency_converter/types"
	"crypto-to-fiat-converter/proto/converter"
)

type service struct {
	cfg *config.Config
	converter.UnimplementedConverterServer
	converter currency_converter.Provider
}

func New(cfg *config.Config, converter currency_converter.Provider) *service {
	return &service{
		cfg:       cfg,
		converter: converter,
	}
}

func (s *service) Convert(ctx context.Context, req *converter.ConvertRequest) (*converter.ConvertResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "request is empty")
	}

	var convertRequest = make([]types.ConvertRequest, len(req.Request))
	for i, val := range req.Request {
		convertRequest[i] = types.ConvertRequest{
			FromToken:  val.FromToken,
			ToCurrency: val.ToCurrency,
			Amount:     val.Amount,
		}
	}

	responseRaw, err := s.converter.Convert(convertRequest)
	if err != nil {
		return nil, err
	}

	var response = make([]*converter.PriceByCurrency, len(responseRaw))
	for i, val := range responseRaw {
		response[i] = &converter.PriceByCurrency{
			FromToken:    val.FromToken,
			ToCurrency:   val.ToCurrency,
			TokensAmount: val.TokenAmount,
			TotalPrice:   val.TotalPrice,
		}
	}

	return &converter.ConvertResponse{
		Prices: response,
	}, nil
}

func (s *service) GetTokenList(ctx context.Context, req *converter.GetTokenListRequest) (*converter.GetTokenListResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "request is empty")
	}

	if req.PageSize > s.cfg.PageSizeLimit {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("page size is limited by %v", s.cfg.PageSizeLimit))
	}

	tokens, nextPageToken, err := s.converter.GetTokenList(req.PageToken, req.PageSize)
	if err != nil {
		return nil, err
	}

	var response = make([]*converter.Token, len(tokens))
	for i, val := range tokens {
		response[i] = &converter.Token{
			ID:     val.ID,
			Name:   val.Name,
			Symbol: val.Symbol,
		}
	}

	return &converter.GetTokenListResponse{
		Tokens:        response,
		NextPageToken: nextPageToken,
	}, nil
}
