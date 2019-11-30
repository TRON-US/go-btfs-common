package grpc

import (
	"context"
	"github.com/tron-us/go-btfs-common/protos/guard"
)

func GuardClient(addr string) *ClientBuilder {
	return builder(addr)
}

type GuardClientBuilder struct {
	ClientBuilder
}

func (g *GuardClientBuilder) WithContext(ctx context.Context, f func(client guard.GuardServiceClient) error) error {
	return g.doWithContext(ctx, f)
}
