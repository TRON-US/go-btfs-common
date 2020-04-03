package grpc

import (
	"context"
	"strings"
	"testing"

	"github.com/tron-us/go-btfs-common/config"
	"github.com/tron-us/go-btfs-common/protos/shared"

	"github.com/tron-us/go-common/v2/constant"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func init() {
	if config.InitTestDB() != nil {
	}
}

func TestRuntimeServer(t *testing.T) {

	pgConMaps := map[string]string{"DB_URL_STATUS": config.DbStatusURL, "DB_URL_GUARD": config.DbGuardURL}
	rdCon := config.RdURL

	runtime := RuntimeServer{DB_URL: pgConMaps, RD_URL: rdCon, serviceName: "hub"}
	report, err := runtime.CheckRuntime(context.Background(), &shared.SignedRuntimeInfoRequest{})
	if err != nil {
		assert.Error(t, err, zap.Error(err))
	}

	assert.True(t, strings.Contains(string(report.ServiceName), "hub"), "service name assigned unsuccessfully")
	assert.True(t, strings.Contains(string(report.DbStatusExtra["DB_URL_GUARD"]), constant.DBConnectionHealthy), "database url name assigned unsuccessfully")
	assert.True(t, strings.Contains(string(report.DbStatusExtra["DB_URL_STATUS"]), constant.DBConnectionHealthy), "database url name assigned unsuccessfully")
	assert.True(t, strings.Contains(string(report.RdStatusExtra), constant.RDConnectionHealthy), "redis urls name assigned unsuccessfully")

}

func TestRuntimeServerDBDNE(t *testing.T) {
	pgConMaps := map[string]string{"DB_URL_STATUS": "", "DB_URL_GUARD": ""}
	rdCon := config.RdURL
	const DBURLDNE = "DB URL does not exist !!"

	runtime := RuntimeServer{DB_URL: pgConMaps, RD_URL: rdCon, serviceName: "hub"}
	report, err := runtime.CheckRuntime(context.Background(), &shared.SignedRuntimeInfoRequest{})
	if err != nil {
		assert.Error(t, err, zap.Error(err))
	}

	assert.True(t, strings.Contains(string(report.ServiceName), "hub"), "service name assigned unsuccessfully")
	assert.True(t, strings.Contains(string(report.DbStatusExtra["DB_URL_STATUS"]), DBURLDNE), "database url name assigned unsuccessfully")
	assert.True(t, strings.Contains(string(report.DbStatusExtra["DB_URL_GUARD"]), DBURLDNE), "database url name assigned unsuccessfully")
	assert.True(t, strings.Contains(string(report.RdStatusExtra), constant.RDConnectionHealthy), "redis urls name assigned unsuccessfully")

}
