package grpc

import (
	"log"
	"net"

	"github.com/tron-us/go-btfs-common/protos/shared"
	pb "github.com/tron-us/go-btfs-common/protos/status"
	"github.com/tron-us/go-common/v2/constant"
	"github.com/tron-us/go-common/v2/middleware"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthPb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	server       *grpc.Server
	healthServer *health.Server
	serverName   string
	lis          net.Listener
	dBURL        string
	rDURL        string
}

func (s *GrpcServer) GrpcStatusServer(port string, dBURL string, rDURL string, statusServer pb.StatusServiceServer) *GrpcServer {
	s.dBURL = dBURL
	s.rDURL = rDURL

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Panic(constant.RPCListenError, zap.Error(err))
	}
	s.lis = lis
	s.CreateServer("status-server").
		CreateHealthServer().
		RegisterServer(statusServer).
		RegisterHealthServer().
		WithReflection().
		AcceptConnection()

	return s
}

func (s *GrpcServer) AcceptConnection() *GrpcServer {
	if err := s.server.Serve(s.lis); err != nil {
		log.Panic(constant.RPCServeError, zap.Error(err))
	}
	return s
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

func (s *GrpcServer) RegisterServer(serverl interface{}) *GrpcServer {
	//register two services under the same server
	pb.RegisterStatusServer(s.server, serverl.(pb.StatusServiceServer))
	pb.RegisterStatusServiceServer(s.server, serverl.(pb.StatusServiceServer))
	shared.RegisterRuntimeServiceServer(s.server, &RuntimeServer{DB_URL: s.dBURL, RD_URL: s.rDURL, serviceName: s.serverName})
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
