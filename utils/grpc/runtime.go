package grpc

import (
	"context"

	sharedpb "github.com/tron-us/go-btfs-common/protos/shared"
)

func RuntimeClient(addr string) *RuntimeClientBuilder {
	return &RuntimeClientBuilder{builder(addr)}
}

type RuntimeClientBuilder struct {
	ClientBuilder
}

func (g *RuntimeClientBuilder) WithContext(ctx context.Context, f func(ctx context.Context,
	client sharedpb.RuntimeServiceClient) error) error {
	return g.doWithContext(ctx, f)
}
