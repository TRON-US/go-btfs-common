package grpc

import (
	"context"
	"time"

	sharedpb "github.com/tron-us/go-btfs-common/protos/shared"
	"github.com/tron-us/go-btfs-common/utils"

	"github.com/tron-us/go-common/v2/db"
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
	DB_URL      map[string]string
	RD_URL      string
	serviceName string
	sharedpb.UnimplementedRuntimeServiceServer
}

var Startime = time.Now()

//implementation of the shared helper function
func (s *RuntimeServer) CheckRuntime(ctx context.Context, req *sharedpb.SignedRuntimeInfoRequest) (*sharedpb.RuntimeInfoReport, error) {
	//get connection object

	connection := db.ConnectionUrls{
		PgURL: s.DB_URL,
		RdURL: s.RD_URL,
	}

	//check runtime in shared
	res, _, err := utils.CheckDBConnection(ctx, req, connection)
	if err != nil {
		return nil, err
	}
	//fill the returned data with server specific info
	if res != nil {
		res.QueueStatusExtra = ""
		res.ChainStatusExtra = ""
		res.CacheStatusExtra = ""
		res.Extra = ""
		res.PeerId = ""
		res.ServiceName = s.serviceName
		res.StartTime = Startime
		res.CurentTime = time.Now()
		res.GitHash = ""
		res.Version = ""
	}

	return res, err
}
