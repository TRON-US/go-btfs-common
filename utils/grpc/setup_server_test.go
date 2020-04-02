package grpc

import (
	"context"
	"github.com/stretchr/testify/assert"
	hubpb "github.com/tron-us/go-btfs-common/protos/hub"
	"github.com/tron-us/go-btfs-common/protos/shared"
	sharedpb "github.com/tron-us/go-btfs-common/protos/shared"
	"github.com/tron-us/go-common/v2/constant"
	"github.com/tron-us/go-common/v2/log"
	"go.uber.org/zap"
	"google.golang.org/grpc/connectivity"
	"strings"
	"testing"
	"time"
)
type serverStruct struct {
	hubpb.UnimplementedHubQueryServiceServer
	sharedpb.UnimplementedRuntimeServiceServer
}
func TestSetupServer(t *testing.T) {

	s := GrpcServer{}
	config := "localhost:50030"
	pgConMaps :=  map[string]string{"DB_URL_STATUS":"postgresql://uchenna:Q@vl321!@localhost:5432/db_status", "DB_URL_GUARD": "postgresql://uchenna:Q@vl321!@localhost:5432/db_guard"}
	rdCon := "redis://uchenna:@127.0.0.1:6379/4"
	quit := make(chan struct{})

	//create server
	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				s.GrpcServer(config, pgConMaps, rdCon, &serverStruct{})
			}
		}
	}()

	time.Sleep(time.Second*3)

	//check setup_server variables
	assert.Equal(t, s.serverName, "hub", "hub server name assigned unsuccessfully")
	assert.NotNil(t,s.healthServer,"health server assigned unsuccessfully")
	assert.NotNil(t,s.server,"server assigned unsuccessfully")
	assert.NotNil(t,s.dBURLs, "database urls assigned unsuccessfully")
	assert.NotNil(t,s.rDURL, "redis urls assigned unsuccessfully")

	tests := []struct {
		in  string
		out connectivity.State
		err bool
	}{
		{in: config, err: true},
	}

	//test server with client, check runtime
	for _, tt := range tests {
		err := RuntimeClient(tt.in).WithContext(context.Background(), func(ctx context.Context,
			client shared.RuntimeServiceClient) error {
			res := requestRuntimeInfo(ctx, client)
			//check runtime information
			assert.True(t, strings.Contains(string(res.DbStatusExtra[0]), constant.DBConnectionHealthy), "database assigned unsuccessfully")
			assert.True(t, strings.Contains(string(res.DbStatusExtra[1]), constant.DBConnectionHealthy), "database assigned unsuccessfully")
			assert.True(t, strings.Contains(string(res.RdStatusExtra), constant.RDConnectionHealthy), "redis assigned unsuccessfully")
			return nil
		})
		if err != nil {
			log.Panic("runtime", zap.Error(err))
		}
	}

	close(quit)
}

func requestRuntimeInfo(ctx context.Context, c shared.RuntimeServiceClient) (*sharedpb.RuntimeInfoReport) {
	req := new(shared.SignedRuntimeInfoRequest)
	res, err := c.CheckRuntime(ctx, req)
	if err != nil {
		log.Panic("client", zap.Error(err))
	}
	return res
}

