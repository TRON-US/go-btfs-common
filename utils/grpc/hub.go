package grpc

import (
	"context"
	"time"

	hubpb "github.com/tron-us/go-btfs-common/protos/hub"
)

func HubQueryClient(addr string) *HubQueryClientBuilder {
	return &HubQueryClientBuilder{builder(addr)}
}

func (b *HubQueryClientBuilder) Timeout(to time.Duration) *HubQueryClientBuilder {
	b.timeout = to
	return b
}

type HubQueryClientBuilder struct {
	ClientBuilder
}

func (g *HubQueryClientBuilder) WithContext(ctx context.Context, f func(ctx context.Context,
	client hubpb.HubQueryServiceClient) error) error {
	return g.doWithContext(ctx, f)
}
