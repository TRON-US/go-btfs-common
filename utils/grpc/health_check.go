package grpc

import (
	"context"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func HealthCheckClient(addr string) *HealthCheckClientBuilder {
	return &HealthCheckClientBuilder{builder(addr)}
}

type HealthCheckClientBuilder struct {
	ClientBuilder
}

func (g *HealthCheckClientBuilder) WithContext(ctx context.Context, f func(ctx context.Context,
	client grpc_health_v1.HealthClient) error) error {
	return g.doWithContext(ctx, f)
}
