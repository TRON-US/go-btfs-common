package grpc

import (
	"log"
	"net"

	"github.com/tron-us/go-btfs-common/protos/escrow"
	"github.com/tron-us/go-btfs-common/protos/guard"
	"github.com/tron-us/go-btfs-common/protos/hub"
	"github.com/tron-us/go-btfs-common/protos/shared"
	"github.com/tron-us/go-btfs-common/protos/status"

	"github.com/tron-us/go-common/v2/constant"
	"github.com/tron-us/go-common/v2/middleware"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
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

func (s *GrpcServer) serverTypeToServerName(server interface{}) {
	switch server.(type) {
	case status.StatusServiceServer:
		s.serverName = "status-server"
	case escrow.EscrowServiceServer:
		s.serverName = "escrow"
	case guard.GuardServiceServer:
		s.serverName = "guard"
	case hub.HubQueryServiceServer:
		s.serverName = "hub"
	default:
		s.serverName = "status-server"
	}
}

func (s *GrpcServer) GrpcServer(port string, dbURL string, rdURL string, server interface{}) *GrpcServer {

	s.serverTypeToServerName(server)

	s.dBURL = dbURL
	s.rDURL = rdURL

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Panic(constant.RPCListenError, zap.Error(err))
	}

	s.lis = lis

	s.CreateServer(s.serverName).
		CreateHealthServer().
		RegisterServer(server).
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
	//create grpc health server
	s.healthServer = health.NewServer()
	return s
}

func (s *GrpcServer) RegisterServer(server interface{}) *GrpcServer {

	switch server.(type) {
	case status.StatusServiceServer:
		status.RegisterStatusServer(s.server, server.(status.StatusServiceServer))
		status.RegisterStatusServiceServer(s.server, server.(status.StatusServiceServer))
	case escrow.EscrowServiceServer:
		escrow.RegisterEscrowServiceServer(s.server, server.(escrow.EscrowServiceServer))
	case guard.GuardServiceServer:
		guard.RegisterGuardServiceServer(s.server, server.(guard.GuardServiceServer))
	case hub.HubQueryServiceServer:
		hub.RegisterHubQueryServiceServer(s.server, server.(hub.HubQueryServiceServer))
	}

	shared.RegisterRuntimeServiceServer(s.server, &RuntimeServer{DB_URL: s.dBURL, RD_URL: s.rDURL, serviceName: s.serverName})

	return s
}

func (s *GrpcServer) RegisterHealthServer() *GrpcServer {
	// Add health server to exchange.
	s.healthServer.SetServingStatus(s.serverName, healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s.server, s.healthServer)
	return s
}

func (s *GrpcServer) WithReflection() *GrpcServer {
	// Reflection api register.
	reflection.Register(s.server)
	return s
}
