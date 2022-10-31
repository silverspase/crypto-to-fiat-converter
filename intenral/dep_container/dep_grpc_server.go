package dep_container

import (
	"fmt"
	"net"

	"crypto-to-fiat-converter/intenral/config"
	"crypto-to-fiat-converter/proto/converter"
	"github.com/sarulabs/di"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const grpcServerDefName = "grpcServerDefName"

// RegisterGrpcServer registers GRPC server dependency.
func RegisterGrpcServer(builder *di.Builder) error {
	return builder.Add(di.Def{
		Name: grpcServerDefName,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get(configDefName).(*config.Config)
			listener, err := net.Listen("tcp", fmt.Sprintf(":%v", cfg.Port))
			if err != nil {
				panic(err)
			}

			server := grpc.NewServer()
			service := ctn.Get(grpcServiceDefName).(converter.ConverterServer)

			converter.RegisterConverterServer(server, service)

			zap.L().Info("starting grpc server", zap.String("port", cfg.Port))

			if err = server.Serve(listener); err != nil {
				zap.L().Fatal("failed to serve", zap.Error(err))
			}

			return server, nil
		},
	})
}

// RunGrpcServer runs GRPC Server.
func (c Container) RunGrpcServer() {
	c.container.Get(grpcServerDefName)
}
