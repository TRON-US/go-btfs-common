package grpc

import (
	"context"
	"github.com/tron-us/go-btfs-common/protos/status"
)

func StatusClient(addr string) *ClientBuilder {
	return builder(addr)
}

type StatusClientBuilder struct {
	ClientBuilder
}

func (g *StatusClientBuilder) WithContext(ctx context.Context, f func(client status.StatusClient) error) error {
	return g.doWithContext(ctx, f)
}
