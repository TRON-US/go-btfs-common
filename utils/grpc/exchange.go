package grpc

import (
	"context"
	exchangepb "github.com/tron-us/go-btfs-common/protos/exchange"
)

func ExchangeClient(addr string) *ExchangeClientBuilder {
	return &ExchangeClientBuilder{builder(addr)}
}

type ExchangeClientBuilder struct {
	ClientBuilder
}

func (g *ExchangeClientBuilder) WithContext(ctx context.Context, f func(ctx context.Context,
	client exchangepb.ExchangeClient) error) error {
	return g.doWithContext(ctx, f)
}
