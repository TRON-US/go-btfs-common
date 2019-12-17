package grpc

import (
	"context"
	"time"

	sharedpb "github.com/tron-us/go-btfs-common/protos/shared"
	"github.com/tron-us/go-btfs-common/utils"

	"github.com/tron-us/go-common/v2/db"
	"github.com/tron-us/go-common/v2/log"

	"go.uber.org/zap"
)

func RuntimeClient(addr string) *RuntimeClientBuilder {
	return &RuntimeClientBuilder{builder(addr)}
}

type RuntimeClientBuilder struct {
	ClientBuilder
}

func (g *RuntimeClientBuilder) WithContext(ctx context.Context, f func(ctx context.Context,
	client sharedpb.RuntimeServiceClient) error) error {
	return g.doWithContext(ctx, f)
}

type RuntimeServer struct {
	DB_URL      string
	RD_URL      string
	serviceName string
	sharedpb.UnimplementedRuntimeServiceServer
}

const (
	runtimeHandlerError = "RuntimeHandler Error!"
)

//implementation of the shared helper function
func (s *RuntimeServer) CheckRuntime(ctx context.Context, req *sharedpb.SignedRuntimeInfoRequest) (*sharedpb.RuntimeInfoReport, error) {
	//get connection object

	connection := db.ConnectionUrls{
		PgURL: s.DB_URL,
		RdURL: s.RD_URL,
	}

	//check runtime in shared
	res, err := utils.CheckRuntime(ctx, req, connection)
	if err != nil {
		log.Error(runtimeHandlerError, zap.Error(err))
	}
	//fill the returned data with status-server specific info
	if res != nil {
		res.QueueStatusExtra = []byte(nil)
		res.ChainStatusExtra = []byte(nil)
		res.CacheStatusExtra = []byte(nil)
		res.Extra = []byte(nil)
		res.PeerId = []byte(nil)
		res.ServiceName = []byte(s.serviceName)
		res.CurentTime = time.Now()
		res.GitHash = []byte(nil)
		res.Version = []byte(nil)
	}

	return res, err
}
