package grpc

import (
	"context"
	tronpb "github.com/tron-us/go-btfs-common/protos/protocol/api"
)

func SolidityClient(addr string) *SolidityClientBuilder {
	return &SolidityClientBuilder{builder(addr)}
}

type SolidityClientBuilder struct {
	ClientBuilder
}

func (g *SolidityClientBuilder) WithContext(ctx context.Context, f func(ctx context.Context,
	client tronpb.WalletSolidityClient) error) error {
	return g.doWithContext(ctx, f)
}
