package utils

import (
	"context"
	"os"
	"testing"
	"time"

	sharedpb "github.com/tron-us/go-btfs-common/protos/shared"
	"github.com/tron-us/go-common/v2/constant"
	"github.com/tron-us/go-common/v2/db"

	"github.com/stretchr/testify/assert"
)

var pgURLString string
var rdURLString string

func init() {
	//get db and redis connection strings
	pgURLString = os.Getenv("TEST_DB_URL")
	rdURLString = os.Getenv("TEST_RD_URL")
}

func TestCheckRuntimeDB(t *testing.T) {
	//setup connection (postgres) object
	var connection = db.ConnectionUrls{
		PgURL: pgURLString,
		RdURL: "",
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	shared := new(sharedpb.SignedRuntimeInfoRequest)
	runtimeInfoReportPass, _ := CheckRuntime(ctx, shared, connection)
	assert.Equal(t, runtimeInfoReportPass.DbStatusExtra, []byte(constant.DBConnectionHealthy), "DB is not running")
	//disable connection string
	connection.PgURL = ""
	runtimeInfoReportFail, _ := CheckRuntime(ctx, shared, connection)
	assert.Equal(t, runtimeInfoReportFail.DbStatusExtra, []byte(nil), "DB connection is still provided, error!")
}
func TestCheckRuntimeRD(t *testing.T) {
	//setup connection (redis) object
	var connection = db.ConnectionUrls{
		PgURL: "",
		RdURL: rdURLString,
	}
	shared := new(sharedpb.SignedRuntimeInfoRequest)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	runtimeInfoReportPass, _ := CheckRuntime(ctx, shared, connection)
	assert.Equal(t, runtimeInfoReportPass.RdStatusExtra, []byte(constant.RDConnectionHealthy), "Redis is not running")
	//disable connection string
	connection.RdURL = ""
	runtimeInfoReportFail, _ := CheckRuntime(ctx, shared, connection)
	assert.Equal(t, runtimeInfoReportFail.RdStatusExtra, []byte(nil), "Redis connection is still provided, error!")
}
