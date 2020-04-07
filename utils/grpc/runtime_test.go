package grpc

import (
	"context"
	"github.com/tron-us/go-common/v2/log"
	"strings"
	"testing"

	"github.com/tron-us/go-btfs-common/config"
	"github.com/tron-us/go-btfs-common/protos/shared"

	"github.com/tron-us/go-common/v2/constant"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func init() {
	if err := config.InitTestDB(); err != nil {
		log.Panic("Cannot init database urls for testing !", zap.Error(err))
	}
}

func TestRuntimeServer(t *testing.T) {

	pgConMaps := map[string]string{"DB_URL_STATUS": config.DbStatusURL, "DB_URL_GUARD": config.DbGuardURL}
	rdCon := config.RdURL

	runtime := RuntimeServer{DB_URL: pgConMaps, RD_URL: rdCon, serviceName: "hub"}
	report, err := runtime.CheckRuntime(context.Background(), &shared.SignedRuntimeInfoRequest{})

	assert.Nil(t, err, zap.Error(err))
	assert.True(t, strings.Contains(report.ServiceName, "hub"), "service name assigned unsuccessfully")
	assert.True(t, strings.Contains(report.DbStatusExtra["DB_URL_GUARD"], constant.DBConnectionHealthy), "database url name assigned unsuccessfully")
	assert.True(t, strings.Contains(report.DbStatusExtra["DB_URL_STATUS"], constant.DBConnectionHealthy), "database url name assigned unsuccessfully")
	assert.True(t, strings.Contains(report.RdStatusExtra, constant.RDConnectionHealthy), "redis urls name assigned unsuccessfully")

}

func TestRuntimeServerDBDNE(t *testing.T) {
	pgConMaps := map[string]string{"DB_URL_STATUS": "", "DB_URL_GUARD": ""}
	rdCon := config.RdURL
	const DBURLDNE = "DB URL does not exist !!"

	runtime := RuntimeServer{DB_URL: pgConMaps, RD_URL: rdCon, serviceName: "hub"}
	report, err := runtime.CheckRuntime(context.Background(), &shared.SignedRuntimeInfoRequest{})

	assert.Nil(t, err, zap.Error(err))
	assert.True(t, strings.Contains(report.ServiceName, "hub"), "service name assigned unsuccessfully")
	assert.True(t, strings.Contains(report.DbStatusExtra["DB_URL_STATUS"], DBURLDNE), "database url name assigned unsuccessfully")
	assert.True(t, strings.Contains(report.DbStatusExtra["DB_URL_GUARD"], DBURLDNE), "database url name assigned unsuccessfully")
	assert.True(t, strings.Contains(report.RdStatusExtra, constant.RDConnectionHealthy), "redis urls name assigned unsuccessfully")

}
