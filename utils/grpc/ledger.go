package grpc

import (
	"context"
	"github.com/tron-us/go-btfs-common/protos/ledger"
)

func LedgerClient(addr string) *LedgerClientBuilder {
	return &LedgerClientBuilder{builder(addr)}
}

type LedgerClientBuilder struct {
	ClientBuilder
}

func (g *LedgerClientBuilder) WithContext(ctx context.Context, f func(ctx context.Context,
	client ledger.ChannelsClient) error) error {
	return g.doWithContext(ctx, f)
}
