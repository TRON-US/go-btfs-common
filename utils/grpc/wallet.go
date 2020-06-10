package grpc

import (
	"context"
	tronpb "github.com/tron-us/go-btfs-common/protos/protocol/api"
)

func WalletClient(addr string) *WalletClientBuilder {
	return &WalletClientBuilder{builder(addr)}
}

type WalletClientBuilder struct {
	ClientBuilder
}

func (g *WalletClientBuilder) WithContext(ctx context.Context, f func(ctx context.Context,
	client tronpb.WalletClient) error) error {
	return g.doWithContext(ctx, f)
}
