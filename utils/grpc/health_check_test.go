package grpc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestHealthCheckClient(t *testing.T) {
	tpl := "https://%s.btfs.io"
	srvs := []struct {
		addr string
		name string
	}{
		{addr: "escrow", name: "escrow"},
		{addr: "guard", name: "guard-interceptor"},
		{addr: "hub", name: "hub-query"},
		{addr: "status", name: "status-server"},
	}
	for _, srv := range srvs {
		client := HealthCheckClient(fmt.Sprintf(tpl, srv.addr))
		client.Timeout(60 * time.Second)
		err := client.WithContext(context.Background(), func(ctx context.Context, client grpc_health_v1.HealthClient) error {
			req := &grpc_health_v1.HealthCheckRequest{Service: srv.name}
			check, err := client.Check(ctx, req)
			if err == nil {
				assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, check.Status)
			}
			return err
		})
		if err != nil {
			t.Fatal(srv, err)
		}
	}
}
