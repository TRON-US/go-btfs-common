package utils

import (
	"context"
	"github.com/tron-us/go-common/env"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sharedpb "github.com/tron-us/go-btfs-common/protos/shared"
	"github.com/tron-us/go-common/constant"
)

var (
	testDbURL = ""
	testDbVal = ""
	testRdURL = ""
	testRdVal = ""
)

func init() {
	//check the environment
	if env.IsDev() {
		//pass db/rd credentials, save them locally
		testDbURL, testDbVal = env.GetEnv("DB_URL")
		testRdURL, testRdVal = env.GetEnv("RD_URL")
	}
}
func initEnvironment(){
	//initialize environment since some tests clear this information
	os.Setenv(testDbURL, testDbVal)
	os.Setenv(testRdURL, testRdVal)
}

func TestCheckRuntimeDBRD(t *testing.T) {
	// setup
	initEnvironment()
	//check environment
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	shared := new(sharedpb.SignedRuntimeInfoRequest)
	runtimeInfoReport, _ := CheckRuntime(ctx, shared)
	assert.Equal(t, runtimeInfoReport.Status, sharedpb.RuntimeInfoReport_RUNNING, "CheckRuntime failed DB status check")
	assert.Equal(t, runtimeInfoReport.DbStatusExtra, []byte(constant.DBConnectionHealthy), "CheckRuntime failed DB status extra check ")
	assert.Equal(t, runtimeInfoReport.RdStatusExtra, []byte(constant.RDConnectionHealthy), "CheckRuntime failed Redis status extra check ")
}

func TestCheckRuntimeDBRDFail(t *testing.T) {
	//setup
	initEnvironment()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	shared := new(sharedpb.SignedRuntimeInfoRequest)
	//unset redis
	os.Unsetenv(testRdURL)
	runtimeInfoReport, _ := CheckRuntime(ctx, shared)
	assert.Equal(t, runtimeInfoReport.Status, sharedpb.RuntimeInfoReport_RUNNING, "CheckRuntime failed DB status check")
	assert.Equal(t, runtimeInfoReport.DbStatusExtra, []byte(constant.DBConnectionHealthy), "CheckRuntime failed DB status extra check ")
	assert.Equal(t, runtimeInfoReport.RdStatusExtra, []byte(nil), "CheckRuntime did not fail Redis status extra check ")
}

func TestCheckRuntimeRDPassDBFail(t *testing.T) {
	//setup
	initEnvironment()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	shared := new(sharedpb.SignedRuntimeInfoRequest)
	//unset db
	os.Unsetenv(testDbURL)
	runtimeInfoReport, _ := CheckRuntime(ctx, shared)
	assert.Equal(t, runtimeInfoReport.Status, sharedpb.RuntimeInfoReport_RUNNING, "CheckRuntime failed DB status check")
	assert.Equal(t, runtimeInfoReport.RdStatusExtra, []byte(constant.RDConnectionHealthy), "CheckRuntime failed Redis status extra check ")
	assert.Equal(t, runtimeInfoReport.DbStatusExtra, []byte(nil), "CheckRuntime did not fail DB status extra check ")
}

func TestCheckRuntimeDBFailRDFail(t *testing.T) {
	//setup
	initEnvironment()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	shared := new(sharedpb.SignedRuntimeInfoRequest)
	//unset both
	os.Unsetenv(testRdURL)
	os.Unsetenv(testDbURL)
	runtimeInfoReport, _ := CheckRuntime(ctx, shared)
	assert.Equal(t, runtimeInfoReport.Status, sharedpb.RuntimeInfoReport_RUNNING, "CheckRuntime failed DB status check")
	assert.Equal(t, runtimeInfoReport.DbStatusExtra, []byte(nil), "CheckRuntime did not fail DB status extra check ")
	assert.Equal(t, runtimeInfoReport.RdStatusExtra, []byte(nil), "CheckRuntime did not fail Redis status extra check ")
}
