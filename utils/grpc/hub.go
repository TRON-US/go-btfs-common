package grpc

import (
	"context"
	"github.com/tron-us/go-btfs-common/protos/hub"
)

func HubQueryClient(addr string) *ClientBuilder {
	return builder(addr)
}

type HubQueryClientBuilder struct {
	ClientBuilder
}

func (g *HubQueryClientBuilder) WithContext(ctx context.Context, f func(client hub.HubQueryClient) error) error {
	return g.doWithContext(ctx, f)
}
