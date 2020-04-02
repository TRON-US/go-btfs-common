package utils

import (
	"context"
	"fmt"
	"github.com/tron-us/go-btfs-common/config"
	"strings"
	"testing"
	"time"

	sharedpb "github.com/tron-us/go-btfs-common/protos/shared"
	"github.com/tron-us/go-common/v2/constant"
	"github.com/tron-us/go-common/v2/db"

	"github.com/stretchr/testify/assert"
)

var dbStatusURL string
var dbGuardURL string
var rdURL string
var foundRDString bool
var foundPgStatusString bool
var foundPgGuardString bool

func TestCheckRuntimeDB(t *testing.T) {
	config.InitDB()
	//setup connection (postgres) object
	fmt.Println(dbStatusURL)
	fmt.Println(dbGuardURL)
	pgConMaps :=  map[string]string{"DB_URL_STATUS":dbStatusURL, "DB_URL_GUARD": dbGuardURL}
	var connection = db.ConnectionUrls{
		PgURL: pgConMaps,
		RdURL: "",
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	shared := new(sharedpb.SignedRuntimeInfoRequest)
	runtimeInfoReportPass, _ := CheckRuntime(ctx, shared, connection)
	assert.Equal(t, runtimeInfoReportPass.DbStatusExtra[0], []byte("DB_URL_STATUS:" +constant.DBConnectionHealthy), "connection successful")
	assert.Equal(t, runtimeInfoReportPass.DbStatusExtra[1], []byte("DB_URL_GUARD:" +constant.DBConnectionHealthy), "connection successful")

	//disable connection string
	pgEmptyConMaps :=  map[string]string{"DB_URL_STATUS":"", "DB_URL_GUARD": ""}
	var emptyConnection = db.ConnectionUrls{
		PgURL: pgEmptyConMaps,
		RdURL: "",
	}
	runtimeInfoReportFail, _ := CheckRuntime(ctx, shared, emptyConnection)
	assert.True(t, strings.Contains(string(runtimeInfoReportFail.DbStatusExtra[0]), "DB URL does not exist !!"), "DB URL does not exist !!")
	assert.True(t, strings.Contains(string(runtimeInfoReportFail.DbStatusExtra[1]), "DB URL does not exist !!"), "DB URL does not exist !!")
}
func TestCheckRuntimeRD(t *testing.T) {
	config.InitDB()
	const RDURLDNE = "RD URL does not exist !!"
	//setup connection (redis) object
	var connection = db.ConnectionUrls{
		PgURL: map[string]string{},
		RdURL: rdURL,
	}
	shared := new(sharedpb.SignedRuntimeInfoRequest)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	runtimeInfoReportPass, _ := CheckRuntime(ctx, shared, connection)
	assert.True(t, strings.Contains(string(runtimeInfoReportPass.RdStatusExtra), constant.RDConnectionHealthy), "Redis is not running")
	//disable connection string
	connection.RdURL = ""
	runtimeInfoReportFail, _ := CheckRuntime(ctx, shared, connection)
	assert.True(t, strings.Contains(string(runtimeInfoReportFail.RdStatusExtra), RDURLDNE) , "Redis connection is still provided, error!")
}
