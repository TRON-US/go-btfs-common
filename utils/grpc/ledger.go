package grpc

import (
	"context"
	ledgerpb "github.com/tron-us/go-btfs-common/protos/ledger"
	"time"
)

func LedgerClient(addr string) *LedgerClientBuilder {
	return &LedgerClientBuilder{builder(addr)}
}

func (b *LedgerClientBuilder) Timeout(to time.Duration) *LedgerClientBuilder {
	b.timeout = to
	return b
}

type LedgerClientBuilder struct {
	ClientBuilder
}

func (g *LedgerClientBuilder) WithContext(ctx context.Context, f func(ctx context.Context,
	client ledgerpb.ChannelsClient) error) error {
	return g.doWithContext(ctx, f)
}
