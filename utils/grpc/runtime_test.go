package grpc

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/tron-us/go-btfs-common/config"
	"github.com/tron-us/go-btfs-common/protos/shared"
	"github.com/tron-us/go-common/v2/constant"
	"github.com/tron-us/go-common/v2/log"
	"go.uber.org/zap"
	"strings"
	"testing"
)

func TestRuntimeServer(t *testing.T) {
	config.InitDB()
	pgConMaps :=  map[string]string{"DB_URL_STATUS":config.DbStatusURL, "DB_URL_GUARD": config.DbGuardURL}
	rdCon := config.RdURL

	runtime := RuntimeServer{DB_URL:pgConMaps, RD_URL:rdCon, serviceName:"hub"}
	report, err := runtime.CheckRuntime(context.Background(), &shared.SignedRuntimeInfoRequest{})
	if err != nil {
		log.Error("Runtime error", zap.Error(err))
	}

	assert.True(t, strings.Contains(string(report.ServiceName), "hub"), "service name assigned unsuccessfully")
	assert.True(t, strings.Contains(string(report.DbStatusExtra[0]),constant.DBConnectionHealthy) , "database url name assigned unsuccessfully")
	assert.True(t, strings.Contains(string(report.DbStatusExtra[1]), constant.DBConnectionHealthy) ,"database url name assigned unsuccessfully")
	assert.True(t, strings.Contains(string(report.RdStatusExtra), constant.RDConnectionHealthy),  "redis urls name assigned unsuccessfully")

}


func TestRuntimeServerDBDNE(t *testing.T) {
	config.InitDB()
	pgConMaps :=  map[string]string{"DB_URL_STATUS":"", "DB_URL_GUARD": ""}
	rdCon := config.RdURL
	const DBURLDNE = "DB URL does not exist !!"

	runtime := RuntimeServer{DB_URL:pgConMaps, RD_URL:rdCon, serviceName:"hub"}
	report, err := runtime.CheckRuntime(context.Background(), &shared.SignedRuntimeInfoRequest{})
	if err != nil {
		log.Error("Runtime error", zap.Error(err))
	}

	assert.True(t, strings.Contains(string(report.ServiceName), "hub"), "service name assigned unsuccessfully")
	assert.True(t, strings.Contains(string(report.DbStatusExtra[0]), "DB_URL_STATUS"+ ":" + DBURLDNE ),"database url name assigned unsuccessfully")
	assert.True(t, strings.Contains(string(report.DbStatusExtra[1]), "DB_URL_GUARD"+ ":" + DBURLDNE) ,"database url name assigned unsuccessfully")
	assert.True(t, strings.Contains(string(report.RdStatusExtra), constant.RDConnectionHealthy),  "redis urls name assigned unsuccessfully")

}
