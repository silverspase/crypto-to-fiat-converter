package integration

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"crypto-to-fiat-converter/proto/converter"
)

const (
	host = "localhost:8080"
)

type ConvertTestSuite struct {
	suite.Suite
}

func TestAddLotToUserFavoritesTestSuite(t *testing.T) {
	suite.Run(t, new(ConvertTestSuite))
}

func (s *ConvertTestSuite) TestConvert_OK() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := converter.NewConverterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Convert(ctx, &converter.ConvertRequest{
		Request: []*converter.SingleConvertRequest{
			{
				FromToken:  "tether",
				ToCurrency: "usd",
				Amount:     1,
			},
		},
	})
	if err != nil {
		log.Fatalf("failed to send grpc request: %v", err)
	}

	usdtPrice := r.Prices[0].TotalPrice
	if usdtPrice > 1.1 || usdtPrice < 0.9 {
		log.Fatalf("unexpected result: stablecoin price should be ~ 1, but got %v", usdtPrice)
	}
}
