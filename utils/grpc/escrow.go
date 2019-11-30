package grpc

import (
	"context"
	"github.com/tron-us/go-btfs-common/protos/escrow"
)

func EscrowClient(addr string) *ClientBuilder {
	return builder(addr)
}

type EscrowClientBuilder struct {
	ClientBuilder
}

func (g *EscrowClientBuilder) WithContext(ctx context.Context, f func(client escrow.EscrowServiceClient) error) error {
	return g.doWithContext(ctx, f)
}
