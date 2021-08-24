package grpc

import (
	"context"
	"github.com/tron-us/go-btfs-common/protos/nft"
)

func NftClient(addr string) *NftClientBuilder {
	return &NftClientBuilder{builder(addr)}
}

type NftClientBuilder struct {
	ClientBuilder
}

func (g *NftClientBuilder) WithContext(ctx context.Context, f func(ctx context.Context,
	client nft.NftServiceClient) error) error {
	return g.doWithContext(ctx, f)
}
