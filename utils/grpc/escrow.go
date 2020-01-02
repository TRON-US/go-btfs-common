package grpc

import (
	"context"
	"time"

	escrowpb "github.com/tron-us/go-btfs-common/protos/escrow"
)

func EscrowClient(addr string) *EscrowClientBuilder {
	return &EscrowClientBuilder{builder(addr)}
}

func (b *EscrowClientBuilder) Timeout(to time.Duration) *EscrowClientBuilder {
	b.timeout = to
	return b
}

type EscrowClientBuilder struct {
	ClientBuilder
}

func (g *EscrowClientBuilder) WithContext(ctx context.Context, f func(ctx context.Context,
	client escrowpb.EscrowServiceClient) error) error {
	return g.doWithContext(ctx, f)
}
