package grpc

import (
	"context"
	"github.com/tron-us/go-btfs-common/protos/guard"
)

func GuardClient(addr string) *GuardClientBuilder {
	return &GuardClientBuilder{builder(addr)}
}

type GuardClientBuilder struct {
	ClientBuilder
}

func (g *GuardClientBuilder) WithContext(ctx context.Context, f func(ctx context.Context, client guard.GuardServiceClient) error) error {
	return g.doWithContext(ctx, f)
}
