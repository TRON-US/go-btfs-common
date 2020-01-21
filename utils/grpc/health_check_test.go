package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc/health/grpc_health_v1"
	"testing"
	"time"
)

func TestHealthCheckClient(t *testing.T) {
	client := HealthCheckClient("https://hub-dev.btfs.io")
	client.Timeout(10 * time.Second)
	err := client.WithContext(context.Background(), func(ctx context.Context, client grpc_health_v1.HealthClient) error {
		req := &grpc_health_v1.HealthCheckRequest{Service: "hub"}
		check, err := client.Check(ctx, req)
		fmt.Print("check", check)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
}
