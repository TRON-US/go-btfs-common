package grpc

import (
	"context"
	"github.com/tron-us/go-btfs-common/protos/status"
)

func StatusClient(addr string) *StatusClientBuilder {
	return &StatusClientBuilder{builder(addr)}
}

type StatusClientBuilder struct {
	ClientBuilder
}

func (g *StatusClientBuilder) WithContext(ctx context.Context, f func(ctx context.Context,
	client status.StatusServiceClient) error) error {
	return g.doWithContext(ctx, f)
}
