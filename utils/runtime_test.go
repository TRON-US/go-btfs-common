package utils

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/tron-us/go-btfs-common/config"
	sharedpb "github.com/tron-us/go-btfs-common/protos/shared"
	"github.com/tron-us/go-common/v2/db"
	"github.com/tron-us/go-common/v2/log"

	"github.com/tron-us/go-common/v2/constant"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func init() {
	if err := config.InitTestDB(); err != nil {
		log.Panic("Cannot init database urls for testing !", zap.Error(err))
	}
}

func TestCheckDBConnection(t *testing.T) {

	//setup connection (postgres) object
	pgConMaps := map[string]string{"DB_URL_STATUS": config.DbStatusURL, "DB_URL_GUARD": config.DbGuardURL}
	var connection = db.ConnectionUrls{
		PgURL: pgConMaps,
		RdURL: "",
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	shared := new(sharedpb.SignedRuntimeInfoRequest)
	runtimeInfoReportPass, _, err := CheckDBConnection(ctx, shared, connection)
	assert.Nil(t, err, zap.Error(err))
	assert.True(t, strings.Contains(runtimeInfoReportPass.DbStatusExtra["DB_URL_STATUS"], constant.DBConnectionHealthy), "connection not successful")
	assert.True(t, strings.Contains(runtimeInfoReportPass.DbStatusExtra["DB_URL_GUARD"], constant.DBConnectionHealthy), "connection not successful")

	//disable connection string
	pgEmptyConMaps := map[string]string{"DB_URL_STATUS": "", "DB_URL_GUARD": ""}
	var emptyConnection = db.ConnectionUrls{
		PgURL: pgEmptyConMaps,
		RdURL: "",
	}
	runtimeInfoReportFail, _, err := CheckDBConnection(ctx, shared, emptyConnection)
	assert.Nil(t, err, zap.Error(err))
	assert.True(t, strings.Contains(runtimeInfoReportFail.DbStatusExtra["DB_URL_STATUS"], "DB URL does not exist !!"), "DB URL does not exist !!")
	assert.True(t, strings.Contains(runtimeInfoReportFail.DbStatusExtra["DB_URL_GUARD"], "DB URL does not exist !!"), "DB URL does not exist !!")
}
func TestCheckDBConnectionRD(t *testing.T) {
	const RDURLDNE = "RD URL does not exist !!"
	//setup connection (redis) object
	var connection = db.ConnectionUrls{
		PgURL: map[string]string{},
		RdURL: config.RdURL,
	}
	shared := new(sharedpb.SignedRuntimeInfoRequest)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	runtimeInfoReportPass, _, err := CheckDBConnection(ctx, shared, connection)
	assert.Nil(t, err, zap.Error(err))
	assert.True(t, strings.Contains(runtimeInfoReportPass.RdStatusExtra, constant.RDConnectionHealthy), "Redis is not running")
	//disable connection string
	connection.RdURL = ""
	runtimeInfoReportFail, _, err := CheckDBConnection(ctx, shared, connection)
	assert.Nil(t, err, zap.Error(err))
	assert.True(t, strings.Contains(runtimeInfoReportFail.RdStatusExtra, RDURLDNE), "Redis connection is still provided, error!")
}
