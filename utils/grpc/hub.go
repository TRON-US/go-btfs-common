package grpc

import (
	"context"
	hubpb "github.com/tron-us/go-btfs-common/protos/hub"
)

// hub-query
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

// hub-parser
func HubParserClient(addr string) *HubParserClientBuilder {
	return &HubParserClientBuilder{builder(addr)}
}

type HubParserClientBuilder struct {
	ClientBuilder
}

func (g *HubParserClientBuilder) WithContext(ctx context.Context, f func(ctx context.Context,
	client hubpb.HubParseServiceClient) error) error {
	return g.doWithContext(ctx, f)
}
