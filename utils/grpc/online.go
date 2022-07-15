package grpc

import (
	"context"
	"github.com/tron-us/go-btfs-common/protos/online"
)

func OnlineClient(addr string) *OnlineClientBuilder {
	return &OnlineClientBuilder{builder(addr)}
}

type OnlineClientBuilder struct {
	ClientBuilder
}

func (g *OnlineClientBuilder) WithContext(ctx context.Context, f func(ctx context.Context,
	client online.OnlineServiceClient) error) error {
	return g.doWithContext(ctx, f)
}
