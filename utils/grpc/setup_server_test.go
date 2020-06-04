package grpc

import (
	"context"
	"github.com/tron-us/go-common/v2/log"
	"strings"
	"testing"
	"time"

	"github.com/tron-us/go-btfs-common/config"
	hubpb "github.com/tron-us/go-btfs-common/protos/hub"
	"github.com/tron-us/go-btfs-common/protos/shared"
	sharedpb "github.com/tron-us/go-btfs-common/protos/shared"

	"github.com/tron-us/go-common/v2/constant"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc/connectivity"
)

type serverStruct struct {
	hubpb.UnimplementedHubQueryServiceServer
	sharedpb.UnimplementedRuntimeServiceServer
}

func init() {
	if err := config.InitTestDB(); err != nil {
		log.Panic("Cannot init database urls for testing !", zap.Error(err))
	}
}
func TestSetupServer(t *testing.T) {

	s := GrpcServer{}

	address := "localhost:50030"
	pgConMaps := map[string]string{"DB_URL_STATUS": config.DbStatusURL, "DB_URL_GUARD": config.DbGuardURL}
	rdCon := config.RdURL

	//create server
	go func() {
		s.GrpcServer(address, pgConMaps, rdCon, &serverStruct{})
	}()

	time.Sleep(time.Second * 3)

	//check setup_server variables
	assert.Equal(t, s.serverName, "hub-query", "hub server name assigned unsuccessfully")
	assert.NotNil(t, s.healthServer, "health server assigned unsuccessfully")
	assert.NotNil(t, s.server, "server assigned unsuccessfully")
	assert.NotNil(t, s.dBURLs, "database urls assigned unsuccessfully")
	assert.NotNil(t, s.rDURL, "redis urls assigned unsuccessfully")

	tests := []struct {
		in  string
		out connectivity.State
		err bool
	}{
		{in: address, err: true},
	}

	//test server with client, check runtime
	for _, tt := range tests {
		err := RuntimeClient(tt.in).WithContext(context.Background(), func(ctx context.Context,
			client shared.RuntimeServiceClient) error {
			res := requestRuntimeInfo(t, ctx, client)
			//check runtime information
			assert.True(t, strings.Contains(res.DbStatusExtra["DB_URL_STATUS"], constant.DBConnectionHealthy), "database assigned unsuccessfully")
			assert.True(t, strings.Contains(res.DbStatusExtra["DB_URL_GUARD"], constant.DBConnectionHealthy), "database assigned unsuccessfully")
			assert.True(t, strings.Contains(res.RdStatusExtra, constant.RDConnectionHealthy), "redis assigned unsuccessfully")
			return nil
		})
		if err != nil {
			assert.Error(t, err, zap.Error(err))
		}
	}
}

func requestRuntimeInfo(t *testing.T, ctx context.Context, c shared.RuntimeServiceClient) *sharedpb.RuntimeInfoReport {
	req := new(shared.SignedRuntimeInfoRequest)
	res, err := c.CheckRuntime(ctx, req)
	if err != nil {
		assert.Error(t, err, zap.Error(err))
	}
	return res
}
