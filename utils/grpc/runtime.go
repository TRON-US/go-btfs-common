package grpc

import (
	"context"
	sharedpb "github.com/tron-us/go-btfs-common/protos/shared"
)

func RuntimeServiceClient(addr string) *RuntimeServiceClientBuilder {
	return &RuntimeServiceClientBuilder{builder(addr)}
}

type RuntimeServiceClientBuilder struct {
	ClientBuilder
}

func (g *RuntimeServiceClientBuilder) WithContext(ctx context.Context, f func(ctx context.Context,
	client sharedpb.RuntimeServiceClient) error) error {
	return g.doWithContext(ctx, f)
}
