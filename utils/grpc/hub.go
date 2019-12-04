package grpc

import (
	"context"

	hubpb "github.com/tron-us/go-btfs-common/protos/hub"
)

func HubQueryClient(addr string) *HubQueryClientBuilder {
	return &HubQueryClientBuilder{builder(addr)}
}

type HubQueryClientBuilder struct {
	ClientBuilder
}

func (g *HubQueryClientBuilder) WithContext(ctx context.Context, f func(ctx context.Context,
	client hubpb.HubQueryServiceClient) error) error {
	return g.doWithContext(ctx, f)
}
