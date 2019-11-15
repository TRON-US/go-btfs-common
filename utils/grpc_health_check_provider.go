package utils

//support GRPC Health Checking Protocol
//refer https://github.com/grpc/grpc/blob/master/doc/health-checking.md
//Used for Kubernets to run grpc-health-probe to verify the health of the service
// Used by call utils.RegisterHealthCheckService(&grpcServer)
import (
	"context"

	"google.golang.org/grpc"
	he "google.golang.org/grpc/health/grpc_health_v1"
)

type GrpcHealthServer struct{}

func (s *GrpcHealthServer) Check(context.Context, *he.HealthCheckRequest) (*he.HealthCheckResponse, error) {
	return &he.HealthCheckResponse{Status: he.HealthCheckResponse_SERVING}, nil
}

func (s *GrpcHealthServer) Watch(*he.HealthCheckRequest, he.Health_WatchServer) error {
	return nil
}

func RegisterHealthCheckService(s *grpc.Server) {
	he.RegisterHealthServer(s, new(GrpcHealthServer))
}
