package grpc

import (
	"context"
	"fmt"
	"github.com/tron-us/go-btfs-common/protos/online"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/tron-us/go-btfs-common/controller"
	"github.com/tron-us/go-btfs-common/protos/escrow"
	"github.com/tron-us/go-btfs-common/protos/guard"
	"github.com/tron-us/go-btfs-common/protos/hub"
	"github.com/tron-us/go-btfs-common/protos/shared"
	"github.com/tron-us/go-btfs-common/protos/status"
	"github.com/tron-us/go-btfs-common/utils"

	"github.com/tron-us/go-common/v2/constant"
	"github.com/tron-us/go-common/v2/db"
	"github.com/tron-us/go-common/v2/db/postgres"
	"github.com/tron-us/go-common/v2/db/redis"
	"github.com/tron-us/go-common/v2/log"
	"github.com/tron-us/go-common/v2/middleware"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	monitor "github.com/hypnoglow/go-pg-monitor"
	"github.com/hypnoglow/go-pg-monitor/gopgv9"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	promMetricsServer = ":8080"
)

func init() {
	if pms := os.Getenv("PROM_METRICS_SERVER"); pms != "" {
		promMetricsServer = pms
	}
}

type GrpcServer struct {
	server       *grpc.Server
	healthServer *health.Server
	serverName   string
	lis          net.Listener
	dBURLs       map[string]string
	rDURL        string
}

func (s *GrpcServer) serverTypeToServerName(server interface{}) {
	switch t := server.(type) {
	case online.OnlineServiceServer:
		s.serverName = "online-server"
	case status.StatusServiceServer:
		s.serverName = "status-server"
	case escrow.EscrowServiceServer:
		s.serverName = "escrow"
	case guard.GuardServiceServer:
		s.serverName = "guard-interceptor"
	case hub.HubQueryServiceServer:
		s.serverName = "hub-query"
	case hub.HubParseServiceServer:
		s.serverName = "hub-parser"
	case *controller.DefaultController:
		s.serverName = fmt.Sprintf("%v", t.ServerName)
	default:
		s.serverName = "unknown"
	}
}

func (s *GrpcServer) GrpcServer(port string, dbURLs map[string]string, rdURL string,
	server interface{}, options ...grpc.ServerOption) (
	map[string]*postgres.TGPGDB, *redis.TGRDDB) {
	s.serverTypeToServerName(server)

	s.dBURLs = dbURLs
	s.rDURL = rdURL

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Panic(constant.RPCListenError, zap.Error(err))
	}

	s.lis = lis

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check connections first
	connection := db.ConnectionUrls{RdURL: rdURL, PgURL: dbURLs}
	req := new(shared.SignedRuntimeInfoRequest)
	_, connObjs, err := utils.CheckDBConnection(ctx, req, connection)
	if err != nil {
		log.Panic("Unable to connect to DB", zap.Error(err))
	}

	// Create server
	s.CreateServer(s.serverName, options...).
		CreateHealthServer().
		RegisterServer(server).
		RegisterHealthServer().
		WithReflection().
		WithGracefulTermDetectAndExec()

	// Add pg metrics to Prometheus
	pgObjs := map[string]*postgres.TGPGDB{}
	var rdObj *redis.TGRDDB
	for _, co := range connObjs {
		if po, ok := co.(*utils.PostgresObj); ok {
			for key, dbObj := range po.DBs {
				mon := monitor.NewMonitor(
					gopgv9.NewObserver(dbObj.DB),
					monitor.NewMetrics(
						monitor.MetricsWithSubsystem("go_pg_"+key),
					),
				)
				log.Info("Starting Postgres DB Prometheus monitoring", zap.String("db", key))
				mon.Open()
				defer mon.Close()
				// Add
				pgObjs[key] = dbObj
			}
		} else if ro, ok := co.(*utils.RedisObj); ok {
			rdObj = ro.DB
		}
	}

	// After all your registrations, make sure all of the Prometheus metrics are initialized.
	grpc_prometheus.Register(s.server)
	grpc_prometheus.EnableHandlingTimeHistogram()

	go func() {
		// Register Prometheus metrics handler.
		http.Handle("/metrics", promhttp.Handler())
		log.Info("Starting Prometheus /metrics",
			zap.String("service", s.serverName), zap.String("address", promMetricsServer))
		err := http.ListenAndServe(promMetricsServer, nil)
		if err != nil {
			log.Panic("Prometheus listening server is shutting down", zap.Error(err))
		}
	}()

	// GRPC entry point
	log.Info("Starting to accept connections", zap.String("service", s.serverName))
	return pgObjs, rdObj
}

func (s *GrpcServer) AcceptConnection() *GrpcServer {
	if err := s.server.Serve(s.lis); err != nil {
		log.Panic(constant.RPCServeError, zap.Error(err))
	}
	return s
}

func (s *GrpcServer) CreateServer(serverName string, options ...grpc.ServerOption) *GrpcServer {
	//create grpc server
	s.serverName = serverName
	options = append(options, middleware.GrpcServerOption)
	s.server = grpc.NewServer(options...)
	return s
}

func (s *GrpcServer) CreateHealthServer() *GrpcServer {
	//create grpc health server
	s.healthServer = health.NewServer()
	return s
}

func (s *GrpcServer) RegisterServer(server interface{}) *GrpcServer {

	switch server.(type) {
	case online.OnlineServiceServer:
		online.RegisterOnlineServiceServer(s.server, server.(online.OnlineServiceServer))
	case status.StatusServiceServer:
		status.RegisterStatusServiceServer(s.server, server.(status.StatusServiceServer))
	case escrow.EscrowServiceServer:
		escrow.RegisterEscrowServiceServer(s.server, server.(escrow.EscrowServiceServer))
	case guard.GuardServiceServer:
		guard.RegisterGuardServiceServer(s.server, server.(guard.GuardServiceServer))
	case hub.HubQueryServiceServer:
		hub.RegisterHubQueryServiceServer(s.server, server.(hub.HubQueryServiceServer))
	case hub.HubParseServiceServer:
		hub.RegisterHubParseServiceServer(s.server, server.(hub.HubParseServiceServer))
	}

	shared.RegisterRuntimeServiceServer(s.server, &RuntimeServer{DB_URL: s.dBURLs, RD_URL: s.rDURL, serviceName: s.serverName})

	log.Info("Registered: " + s.serverName + " and runtime server!")

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

func (s *GrpcServer) WithGracefulTermDetectAndExec() *GrpcServer {
	//spin another routine to continue execution
	go func() {
		GracefulTerminateDetect()
		GracefulTerminateExec(s.server)
	}()
	return s
}
