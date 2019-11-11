package utils

import (
	"context"
	"os"
	"testing"
	"time"

	sharedpb "github.com/tron-us/go-btfs-common/protos/shared"
	"github.com/tron-us/go-common/constant"

	"github.com/stretchr/testify/assert"
)

func TestCheckRuntimeDBRD(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	shared := new(sharedpb.RuntimeInfoRequest)
	os.Setenv("SS_DB_URL", "postgres://uchenna:Q@vl321!@localhost:5432/db_status")
	os.Setenv("SS_RD_URL", "redis://uchenna:@127.0.0.1:6379/4?pool=25&process=2")
	runtimeInfoReport, _ := CheckRuntime(ctx, shared)
	assert.Equal(t, runtimeInfoReport.Status, sharedpb.RuntimeInfoReport_RUNNING, "CheckRuntime failed DB status check")
	assert.Equal(t, runtimeInfoReport.DbStatusExtra, []byte(constant.DBConnectionHealthy), "CheckRuntime failed DB status extra check ")
	assert.Equal(t, runtimeInfoReport.RdStatusExtra, []byte(constant.RDConnectionHealthy), "CheckRuntime failed Redis status extra check ")
}

func TestCheckRuntimeDBRDFail(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	shared := new(sharedpb.RuntimeInfoRequest)
	os.Unsetenv("SS_DB_URL")
	os.Unsetenv("SS_RD_URL")
	runtimeInfoReport, _ := CheckRuntime(ctx, shared)
	assert.Equal(t, runtimeInfoReport.Status, sharedpb.RuntimeInfoReport_RUNNING, "CheckRuntime failed DB status check")
	assert.Equal(t, runtimeInfoReport.DbStatusExtra, []byte(nil), "CheckRuntime did not fail DB status extra check ")
	assert.Equal(t, runtimeInfoReport.RdStatusExtra, []byte(nil), "CheckRuntime did not fail Redis status extra check ")
}

func TestCheckRuntimeRDPassDBFail(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	shared := new(sharedpb.RuntimeInfoRequest)
	os.Setenv("SS_RD_URL", "redis://uchenna:@127.0.0.1:6379/4?pool=25&process=2")
	os.Unsetenv("SS_DB_URL")
	runtimeInfoReport, _ := CheckRuntime(ctx, shared)
	assert.Equal(t, runtimeInfoReport.Status, sharedpb.RuntimeInfoReport_RUNNING, "CheckRuntime failed DB status check")
	assert.Equal(t, runtimeInfoReport.RdStatusExtra, []byte(constant.RDConnectionHealthy), "CheckRuntime failed Redis status extra check ")
	assert.Equal(t, runtimeInfoReport.DbStatusExtra, []byte(nil), "CheckRuntime did not fail DB status extra check ")
}

func TestCheckRuntimeDBPassRDFail(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	shared := new(sharedpb.RuntimeInfoRequest)
	os.Setenv("SS_DB_URL", "postgres://uchenna:Q@vl321!@localhost:5432/db_status")
	os.Unsetenv("SS_RD_URL")
	runtimeInfoReport, _ := CheckRuntime(ctx, shared)
	assert.Equal(t, runtimeInfoReport.Status, sharedpb.RuntimeInfoReport_RUNNING, "CheckRuntime failed DB status check")
	assert.Equal(t, runtimeInfoReport.DbStatusExtra, []byte(constant.DBConnectionHealthy), "CheckRuntime failed DB status extra check ")
	assert.Equal(t, runtimeInfoReport.RdStatusExtra, []byte(nil), "CheckRuntime did not fail Redis status extra check ")
}
