package dep_container

import (
	"crypto-to-fiat-converter/intenral/config"
	"crypto-to-fiat-converter/intenral/service/currency_converter"
	"crypto-to-fiat-converter/intenral/service/grpc"
	"github.com/sarulabs/di"
)

const grpcServiceDefName = "grpcServiceDefName"

// RegisterGrpcService registers GRPC service dependency that goes as a part of GRPC server
func RegisterGrpcService(builder *di.Builder) error {
	return builder.Add(di.Def{
		Name: grpcServiceDefName,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get(configDefName).(*config.Config)
			currencyConverter := ctn.Get(currencyConverterServiceDefName).(currency_converter.Provider)

			return grpc_service.New(cfg, currencyConverter), nil
		},
	})
}
