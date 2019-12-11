package grpc

import (
	"github.com/tron-us/go-btfs-common/protos/shared"
	pb "github.com/tron-us/go-btfs-common/protos/status"
	"github.com/tron-us/go-common/v2/middleware"
	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthPb "google.golang.org/grpc/health/grpc_health_v1"
)

type RuntimeServer struct {
	shared.UnimplementedRuntimeServiceServer
}
type StatusServer struct {
	// for backward-compatibility
	//pb.UnimplementedStatusServer
	pb.UnimplementedStatusServiceServer
}
type GrpcServer struct {
	server       *grpc.Server
	healthServer *health.Server
	serverName   string
}

func (s *GrpcServer) CreateServer(serverName string) *GrpcServer {
	//create grpc server
	s.serverName = serverName
	s.server = grpc.NewServer(middleware.GrpcServerOption)
	return s
}
func (s *GrpcServer) CreateHealthServer() *GrpcServer {
	//create grpc heath server
	s.healthServer = health.NewServer()
	return s
}
func (s *GrpcServer) RegisterServer() *GrpcServer {
	//register two services under the same server
	pb.RegisterStatusServer(s.server, &StatusServer{})
	pb.RegisterStatusServiceServer(s.server, &StatusServer{})
	shared.RegisterRuntimeServiceServer(s.server, &RuntimeServer{})
	return s
}

func (s *GrpcServer) RegisterHealthServer() *GrpcServer {
	// Add health server to exchange.
	s.healthServer.SetServingStatus(s.serverName, healthPb.HealthCheckResponse_SERVING)
	healthPb.RegisterHealthServer(s.server, s.healthServer)
	return s
}

func (s *GrpcServer) WithReflection() *GrpcServer {
	// Reflection api register.
	reflection.Register(s.server)
	return s
}

func (s *GrpcServer) GetServer() *grpc.Server {
	return s.server
}
