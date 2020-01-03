package grpc

import (
	"context"
	escrowpb "github.com/tron-us/go-btfs-common/protos/escrow"
)

func EscrowClient(addr string) *EscrowClientBuilder {
	return &EscrowClientBuilder{builder(addr)}
}

type EscrowClientBuilder struct {
	ClientBuilder
}

func (g *EscrowClientBuilder) WithContext(ctx context.Context, f func(ctx context.Context,
	client escrowpb.EscrowServiceClient) error) error {
	return g.doWithContext(ctx, f)
}
